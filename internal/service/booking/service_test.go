package booking_test

import (
	"testing"
	"time"

	"gotus/internal/mocks"
	"gotus/internal/model/book"
	"gotus/internal/model/reservation"
	"gotus/internal/model/user"
	"gotus/internal/service/booking"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func newBookingServiceWithMocks(ctrl *gomock.Controller) (*booking.Service, *mocks.MockUserRepository, *mocks.MockBookRepository, *mocks.MockBookInstanceRepository, *mocks.MockReservationRepository, *mocks.MockLogger) {
	userRepo := mocks.NewMockUserRepository(ctrl)
	bookRepo := mocks.NewMockBookRepository(ctrl)
	instanceRepo := mocks.NewMockBookInstanceRepository(ctrl)
	resRepo := mocks.NewMockReservationRepository(ctrl)
	logger := mocks.NewMockLogger(ctrl)
	service := booking.NewBookingService(userRepo, bookRepo, instanceRepo, resRepo, logger)
	logger.EXPECT().Log(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	return &service, userRepo, bookRepo, instanceRepo, resRepo, logger
}

func TestCreateBooking_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	svc, userRepo, _, _, _, _ := newBookingServiceWithMocks(ctrl)

	userRepo.EXPECT().FindUserById(1).Return(nil, false)
	res, err := (*svc).CreateBooking(1, "978-0-00-000000-0")
	assert.Nil(t, res)
	assert.ErrorIs(t, err, booking.ErrUserNotFound)
}

func TestCreateBooking_BookNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	svc, userRepo, bookRepo, _, _, _ := newBookingServiceWithMocks(ctrl)

	userRepo.EXPECT().FindUserById(1).Return(&user.User{}, true)
	bookRepo.EXPECT().FindBookByISBN("isbn").Return(nil, false)

	res, err := (*svc).CreateBooking(1, "isbn")
	assert.Nil(t, res)
	assert.ErrorIs(t, err, booking.ErrBookNotFound)
}

func TestCreateBooking_NoAvailableInstances(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	svc, userRepo, bookRepo, instanceRepo, resRepo, _ := newBookingServiceWithMocks(ctrl)

	userRepo.EXPECT().FindUserById(1).Return(&user.User{}, true)
	bookRepo.EXPECT().FindBookByISBN("isbn").Return(&book.Book{}, true)
	instanceRepo.EXPECT().GetBookInstancesByISBN("isbn").Return([]*book.BookInstance{
		book.NewBookInstance(101, "isbn"),
	}, 1)
	resRepo.EXPECT().HasActiveReservation(101).Return(true)

	res, err := (*svc).CreateBooking(1, "isbn")
	assert.Nil(t, res)
	assert.ErrorIs(t, err, booking.ErrNoAvailableInstances)
}

func TestExtendBooking_InvalidStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc, _, _, _, resRepo, _ := newBookingServiceWithMocks(ctrl)

	tests := []struct {
		name     string
		mock     func()
		statusID int
		wantOk   bool
		wantErr  error
	}{
		{
			name:     "StatusCancelled_ShouldFail",
			statusID: int(reservation.StatusCancelled),
			mock: func() {
				resRepo.EXPECT().
					FindReservationById(1).
					Return(reservation.NewReservation(1, 1, 1, reservation.StatusCancelled, time.Now(), time.Now().Add(3*24*time.Hour)), true)
			},
			wantOk:  false,
			wantErr: booking.ErrInvalidStatus,
		},
		{
			name:     "StatusEnded_ShouldFail",
			statusID: int(reservation.StatusExtended),
			mock: func() {
				resRepo.EXPECT().
					FindReservationById(1).
					Return(reservation.NewReservation(1, 1, 1, reservation.StatusEnded, time.Now(), time.Now().Add(3*24*time.Hour)), true)
			},
			wantOk:  false,
			wantErr: booking.ErrInvalidStatus,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			ok, err := (*svc).ExtendBooking(1, 3)

			require.Equal(t, tt.wantOk, ok)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestBookingService_ExtendBooking_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc, _, _, _, resRepo, _ := newBookingServiceWithMocks(ctrl)

	now := time.Now()

	tests := []struct {
		name     string
		statusID int
		mock     func()
		wantOk   bool
		wantErr  error
	}{
		{
			name:     "StatusBooked_ShouldSucceed",
			statusID: int(reservation.StatusBooked),
			mock: func() {
				res := reservation.NewReservation(1, 1, 1, reservation.StatusExtended, now.AddDate(0, 0, -10), now.AddDate(0, 0, 7))

				resRepo.EXPECT().
					FindReservationById(1).
					Return(res, true)

				resRepo.EXPECT().
					UpdateReservationById(1, gomock.Any()).
					Return(true)
			},
			wantOk:  true,
			wantErr: nil,
		},
		{
			name:     "StatusExtended_ShouldSucceed",
			statusID: int(reservation.StatusEnded),
			mock: func() {
				res := reservation.NewReservation(1, 1, 1, reservation.StatusExtended, now.AddDate(0, 0, -10), now.AddDate(0, 0, 7))

				resRepo.EXPECT().
					FindReservationById(1).
					Return(res, true)

				resRepo.EXPECT().
					UpdateReservationById(1, gomock.Any()).
					Return(true)
			},
			wantOk:  true,
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			ok, err := (*svc).ExtendBooking(1, 3)

			require.Equal(t, tt.wantOk, ok)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestCancelBooking_NotToday(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	svc, _, _, _, resRepo, _ := newBookingServiceWithMocks(ctrl)

	res := reservation.NewReservation(1, 1, 1, int(reservation.StatusBooked), time.Now().AddDate(0, 0, -1), time.Now())
	resRepo.EXPECT().FindReservationById(1).Return(res, true)

	ok, err := (*svc).CancelBooking(1)
	assert.False(t, ok)
	assert.ErrorIs(t, err, booking.ErrCancelDateMismatch)
}

func TestEndBooking_Today(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	svc, _, _, _, resRepo, _ := newBookingServiceWithMocks(ctrl)

	now := time.Now()
	res := reservation.NewReservation(1, 1, 1, int(reservation.StatusBooked), now, now.AddDate(0, 0, 1))
	resRepo.EXPECT().FindReservationById(1).Return(res, true)

	ok, err := (*svc).EndBooking(1)
	assert.False(t, ok)
	assert.ErrorIs(t, err, booking.ErrEndDateMismatch)
}

func TestCreateBooking_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service, mockUserRepo, mockBookRepo, mockInstanceRepo, mockReservationRepo, _ := newBookingServiceWithMocks(ctrl)

	userID := 1
	isbn := "978-1-56619-909-4"
	bookInstanceID := 42

	mockUserRepo.EXPECT().FindUserById(userID).Return(&user.User{}, true)
	mockBookRepo.EXPECT().FindBookByISBN(isbn).Return(&book.Book{}, true)
	mockInstanceRepo.EXPECT().GetBookInstancesByISBN(isbn).Return([]*book.BookInstance{
		book.NewBookInstance(bookInstanceID, isbn),
	}, 1)
	mockReservationRepo.EXPECT().HasActiveReservation(bookInstanceID).Return(false)
	mockReservationRepo.EXPECT().StoreReservation(gomock.Any())

	res, err := (*service).CreateBooking(userID, isbn)
	assert.NoError(t, err)
	assert.Equal(t, userID, res.UserID)
	assert.Equal(t, bookInstanceID, res.BookInstanceID)
	assert.Equal(t, int(reservation.StatusBooked), res.ReservationStatusID)
}

func TestExtendBooking_TooLong(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service, _, _, _, _, _ := newBookingServiceWithMocks(ctrl)

	success, err := (*service).ExtendBooking(1, 10)
	assert.False(t, success)
	assert.ErrorIs(t, err, booking.ErrExtensionTooLong)
}

func TestCancelBooking_InvalidDate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service, _, _, _, mockReservationRepo, _ := newBookingServiceWithMocks(ctrl)

	res := reservation.NewReservation(1, 1, 1, int(reservation.StatusBooked), time.Now().AddDate(0, 0, -1), time.Now().AddDate(0, 0, 6))
	mockReservationRepo.EXPECT().FindReservationById(1).Return(res, true)

	success, err := (*service).CancelBooking(1)
	assert.False(t, success)
	assert.ErrorIs(t, err, booking.ErrCancelDateMismatch)
}

func TestEndBooking_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc, _, _, _, mockReservationRepo, _ := newBookingServiceWithMocks(ctrl)

	now := time.Now()

	tests := []struct {
		name     string
		statusID int
		mock     func()
		wantOk   bool
		wantErr  error
	}{
		{
			name:     "StatusBooked_ShouldSucceed",
			statusID: int(reservation.StatusBooked),
			mock: func() {
				res := reservation.NewReservation(1, 1, 1, reservation.StatusExtended, now.AddDate(0, 0, -10), now.AddDate(0, 0, 7))

				mockReservationRepo.EXPECT().
					FindReservationById(1).
					Return(res, true)

				mockReservationRepo.EXPECT().
					UpdateReservationById(1, gomock.Any()).
					Return(true)
			},
			wantOk:  true,
			wantErr: nil,
		},
		{
			name:     "StatusExtended_ShouldSucceed",
			statusID: int(reservation.StatusEnded),
			mock: func() {
				res := reservation.NewReservation(1, 1, 1, reservation.StatusExtended, now.AddDate(0, 0, -10), now.AddDate(0, 0, 7))

				mockReservationRepo.EXPECT().
					FindReservationById(1).
					Return(res, true)

				mockReservationRepo.EXPECT().
					UpdateReservationById(1, gomock.Any()).
					Return(true)
			},
			wantOk:  true,
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			ok, err := (*svc).EndBooking(1)

			require.Equal(t, tt.wantOk, ok)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestCancelBooking_ReservationNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	svc, _, _, _, repo, _ := newBookingServiceWithMocks(ctrl)

	repo.EXPECT().FindReservationById(1).Return(nil, false)

	success, err := (*svc).CancelBooking(1)
	assert.False(t, success)
	assert.ErrorIs(t, err, booking.ErrReservationNotFound)
}

func TestCancelBooking__Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, _, _, _, mockReservationRepo, _ := newBookingServiceWithMocks(ctrl)

	res := reservation.NewReservation(1, 1, 1, int(reservation.StatusEnded), time.Now(), time.Now())
	mockReservationRepo.EXPECT().FindReservationById(1).Return(res, true)
	mockReservationRepo.EXPECT().UpdateReservationById(1, gomock.Any()).Return(true)

	success, err := (*service).CancelBooking(1)
	assert.True(t, success)
	assert.NoError(t, err)
}

func TestEndBooking__InvalidStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc, _, _, _, mockReservationRepo, _ := newBookingServiceWithMocks(ctrl)

	tests := []struct {
		name     string
		mock     func()
		statusID int
		wantOk   bool
		wantErr  error
	}{
		{
			name:     "StatusCancelled_ShouldFail",
			statusID: int(reservation.StatusCancelled),
			mock: func() {
				mockReservationRepo.EXPECT().
					FindReservationById(1).
					Return(reservation.NewReservation(1, 1, 1, reservation.StatusCancelled, time.Now().Add(-3*24*time.Hour), time.Now().Add(3*24*time.Hour)), true)
			},
			wantOk:  false,
			wantErr: booking.ErrInvalidStatus,
		},
		{
			name:     "StatusEnded_ShouldFail",
			statusID: int(reservation.StatusEnded),
			mock: func() {
				mockReservationRepo.EXPECT().
					FindReservationById(1).
					Return(reservation.NewReservation(1, 1, 1, reservation.StatusEnded, time.Now().Add(-3*24*time.Hour), time.Now().Add(3*24*time.Hour)), true)
			},
			wantOk:  false,
			wantErr: booking.ErrInvalidStatus,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			ok, err := (*svc).EndBooking(1)

			require.Equal(t, tt.wantOk, ok)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}
