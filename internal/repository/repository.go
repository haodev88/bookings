package repository

import (
	"github.com/haodev88/bookings/internal/models"
	"time"
)

type DatabaseRepo interface {
	AllUsers() bool
	InsertReservation(res models.Reservation) (int, error)
	InsertRoomRestriction(res models.RoomRestriction) error
	SearchAvailabilityByDatesByRoomId(start, end time.Time, roomID int) (bool, error)
	SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error)
	GetRoomById(id int) (models.Room, error)

	GetUserByID(id int) (models.User, error)
	UpdateUser(u models.User) error
	Authenticate(email, testPassword string) (int, string, error)
	AllReservations() ([]models.Reservation, error)
	AllNewReservations() ([]models.Reservation, error)
	GetReservationByID(id int) (models.Reservation, error)
	UpdateReservation(u models.Reservation) error
	DeleteReservation(id int) error
	UpdateProcessedForReservation(id int, processed int) error
	AllRooms() ([]models.Room, error)
	GetRestrictionsForRoomByDate(roomID int, start, end time.Time) ([]models.RoomRestriction, error)
	InsertBlockForRooms(id int, startDate time.Time) error
	DeleteBlockById(id int) error
}
