package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	pb "github.com/MarcGrol/go-training/solutions/hospital/appointments/appointmentapi"
	"github.com/MarcGrol/go-training/solutions/hospital/notifications/notificationapi"
	"github.com/MarcGrol/go-training/solutions/hospital/patients/patientinfoapi"
)

type server struct {
	listener net.Listener
	pb.UnimplementedAppointmentExternalServer
	pb.UnimplementedAppointmentInternalServer

	appointmentStore   AppointmentStore
	patientInfoClient  patientinfoapi.PatientInfoClient
	notificationClient notificationapi.NotificationClient
}

func newServer(store AppointmentStore, patientInfoClient patientinfoapi.PatientInfoClient, notificationClient notificationapi.NotificationClient) *server {
	return &server{
		appointmentStore:   store,
		patientInfoClient:  patientInfoClient,
		notificationClient: notificationClient,
	}
}

func (s *server) GRPCListenBlocking(port string) error {
	var err error
	s.listener, err = net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf("failed to listen at port %s: %v", port, err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAppointmentExternalServer(grpcServer, s)
	pb.RegisterAppointmentInternalServer(grpcServer, s)

	log.Printf("GRPPC server starts listening at port %s...", port)
	err = grpcServer.Serve(s.listener)
	if err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}
	return nil
}

func (s *server) GetAppointmentsOnUser(c context.Context, in *pb.GetAppointmentsOnUserRequest) (*pb.GetAppointmentsReply, error) {
	// Validate input
	// TODO

	// Perform lookup
	internalAppointments, err := s.appointmentStore.GetAppointmentsOnUserUid(in.UserUid)
	if err != nil {
		return &pb.GetAppointmentsReply{
			Error: &pb.Error{
				Code:    500,
				Message: "Technical error fetching appointments on user",
				Details: err.Error(),
			},
		}, nil
	}

	return returnAppointmentList(internalAppointments), nil
}

func (s *server) GetAppointmentsOnStatus(c context.Context, in *pb.GetAppointmentsOnStatusRequest) (*pb.GetAppointmentsReply, error) {
	// Validate input
	// TODO

	// Perform lookup
	internalAppointments, err := s.appointmentStore.GetAppointmentsOnStatus(AppointmentStatus(in.Status))
	if err != nil {
		return &pb.GetAppointmentsReply{
			Error: &pb.Error{
				Code:    500,
				Message: "Technical error fetching appointments on user",
				Details: err.Error(),
			},
		}, nil
	}

	return returnAppointmentList(internalAppointments), nil
}

func (s *server) RequestAppointment(c context.Context, in *pb.RequestAppointmentRequest) (*pb.AppointmentReply, error) {
	// Validate input
	// TODO

	// Check if patient exists
	getPatientResponse, err := s.patientInfoClient.GetPatientOnUid(c, &patientinfoapi.GetPatientOnUidRequest{PatientUid: in.Appointment.UserUid})
	if err != nil {
		log.Printf("Error getting patient %s on uid: %s", in.Appointment.UserUid, err)
		return &pb.AppointmentReply{
			Error: &pb.Error{
				Code:    500,
				Message: "Technical error finding patient on uid",
				Details: err.Error(),
			},
		}, nil
	}
	if getPatientResponse.Error != nil {
		log.Printf("Error getting patient %s on uid: %+v", in.Appointment.UserUid, getPatientResponse.Error)
		return &pb.AppointmentReply{
			Error: &pb.Error{
				Code:    getPatientResponse.Error.Code,
				Message: fmt.Sprintf("Error getting patient on uid:%s", getPatientResponse.Error.Message),
				Details: getPatientResponse.Error.Details,
			},
		}, nil
	}

	// Adjust datastore
	appointmentCreated, err := s.appointmentStore.PutAppointment(convertIntoInternal(*in.Appointment))
	if err != nil {
		log.Printf("Error creating appointment %+v: %s", *in.Appointment, err)
		return &pb.AppointmentReply{
			Error: &pb.Error{
				Code:    500,
				Message: "Technical error creating appointment",
				Details: err.Error(),
			},
		}, nil
	}
	log.Printf("Persisted new appointment")

	return returnSingleAppointment(appointmentCreated), nil
}

func (s *server) ModifyAppointmentStatus(c context.Context, in *pb.ModifyAppointmentStatusRequest) (*pb.AppointmentReply, error) {
	// Validate input
	// TODO

	// Perform lookup
	internalAppointment, found, err := s.appointmentStore.GetAppointmentOnUid(in.AppointmentUid)
	if err != nil {
		log.Printf("Error getting appointment %s: %s", in.AppointmentUid, err)
		return &pb.AppointmentReply{
			Error: &pb.Error{
				Code:    500,
				Message: "Error getting appointment on uid",
				Details: err.Error(),
			},
		}, nil
	}
	if !found {
		log.Printf("Appointment %s not found", in.AppointmentUid)
		return &pb.AppointmentReply{
			Error: &pb.Error{
				Code:    404,
				Message: "Appointment with uid not found",
			},
		}, nil
	}
	log.Printf("Got appointment:%+v", internalAppointment)

	// Fetch patient details
	resp, err := s.patientInfoClient.GetPatientOnUid(c, &patientinfoapi.GetPatientOnUidRequest{PatientUid: internalAppointment.UserUID})
	if err != nil {
		log.Printf("Error getting patient %s on uid: %s", internalAppointment.UserUID, err)
		return &pb.AppointmentReply{
			Error: &pb.Error{
				Code:    500,
				Message: "Technical error finding patient on uid",
				Details: err.Error(),
			},
		}, nil
	}
	if resp.Error != nil {
		log.Printf("Error getting patient %s on uid: %+v", internalAppointment.UserUID, resp.Error)
		return &pb.AppointmentReply{
			Error: &pb.Error{
				Code:    resp.Error.Code,
				Message: fmt.Sprintf("Error getting patient on uid:%s", resp.Error.Message),
				Details: resp.Error.Details,
			},
		}, nil
	}

	log.Printf("Got patient:%+v", resp.Patient)

	// Send out email
	sendEmailResponse, err := s.notificationClient.SendEmail(c, &notificationapi.SendEmailRequest{
		Email: &notificationapi.EmailMessage{
			RecipientEmailAddress: resp.Patient.PhoneNumber,
			Subject:               "Appointment confirmed",              // TODO use template
			Body:                  "Appointment confirmed with details", // TODO use template
		},
	})
	if err != nil {
		log.Printf("Error sending sms: %s", err)
		return &pb.AppointmentReply{
			Error: &pb.Error{
				Code:    500,
				Message: "Technical error sending sms",
				Details: err.Error(),
			},
		}, nil
	}
	if sendEmailResponse.Error != nil {
		log.Printf("Error sending sms: %+v", sendEmailResponse.Error)
		return &pb.AppointmentReply{
			Error: &pb.Error{
				Code:    sendEmailResponse.Error.Code,
				Message: sendEmailResponse.Error.Message,
				Details: sendEmailResponse.Error.Details,
			},
		}, nil
	}
	log.Printf("Send email:")

	// Send out sms
	sendSmsResponse, err := s.notificationClient.SendSms(c, &notificationapi.SendSmsRequest{
		Sms: &notificationapi.SmsMessage{
			RecipientPhoneNumber: resp.Patient.PhoneNumber,
			Body:                 "Appointment confirmed", // TODO use template
		},
	})
	if err != nil {
		log.Printf("Error sending sms: %s", err)
		return &pb.AppointmentReply{
			Error: &pb.Error{
				Code:    500,
				Message: "Technical error sending sms",
				Details: err.Error(),
			},
		}, nil
	}
	if sendSmsResponse.Error != nil {
		log.Printf("Error sending sms: %+v", sendSmsResponse.Error)
		return &pb.AppointmentReply{
			Error: &pb.Error{
				Code:    sendSmsResponse.Error.Code,
				Message: sendSmsResponse.Error.Message,
				Details: sendSmsResponse.Error.Details,
			},
		}, nil
	}
	log.Printf("Send sms:")

	// Adjust datastore
	internalAppointment.Status = AppointmentStatusConfirmed
	appointmentAdjusted, err := s.appointmentStore.PutAppointment(internalAppointment)
	if err != nil {
		return &pb.AppointmentReply{
			Error: &pb.Error{
				Code:    500,
				Message: "Error persisting modified appointment",
				Details: resp.Error.Details,
			},
		}, nil
	}
	log.Printf("Persisted adjusted appointment")

	return returnSingleAppointment(appointmentAdjusted), nil
}

func returnAppointmentList(internalAppointments []Appointment) *pb.GetAppointmentsReply {
	externalAppointments := []*pb.Appointment{}
	for _, a := range internalAppointments {
		externalAppointments = append(externalAppointments, convertIntoExternal(a))
	}
	return &pb.GetAppointmentsReply{
		Appointments: externalAppointments,
	}
}

func returnSingleAppointment(internalAppointmnent Appointment) *pb.AppointmentReply {
	return &pb.AppointmentReply{
		Appointment: convertIntoExternal(internalAppointmnent),
	}
}

func convertIntoExternal(a Appointment) *pb.Appointment {
	return &pb.Appointment{
		AppointmentUid: a.AppointmentUID,
		UserUid:        a.UserUID,
		DateTime:       a.DateTime,
		Location:       a.Location,
		Topic:          a.Topic,
		Status:         pb.AppointmentStatus(a.Status),
	}
}

func convertIntoInternal(a pb.Appointment) Appointment {
	return Appointment{
		AppointmentUID: a.AppointmentUid,
		UserUID:        a.UserUid,
		DateTime:       a.DateTime,
		Location:       a.Location,
		Topic:          a.Topic,
		Status:         AppointmentStatus(a.Status),
	}
}
