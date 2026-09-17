package appointment

import (
	"errors"
	"testing"
	"time"
	"uuid"
)

func mustNewAppointment(t *testing.T, serviceIDs []uuid.UUID) *Appointment {
	t.Helper()

	start := time.Date(2026, time.September, 18, 10, 0, 0, 0, time.UTC)
	end := start.Add(time.Hour)

	return mustNewAppointmentWithTimes(t, serviceIDs, start, end)
}

func mustNewAppointmentWithTimes(
	t *testing.T,
	serviceIDs []uuid.UUID,
	start, end time.Time,
) *Appointment {
	t.Helper()

	userID := uuid.NewV7()
	masterID := uuid.NewV7()

	appt, err := NewAppointment(userID, masterID, serviceIDs, start, end)
	if err != nil {
		t.Fatalf("NewAppointment() unexpected error: %v", err)
	}

	if appt == nil {
		t.Fatal("NewAppointment() returned nil")
	}

	return appt
}

func containsUUID(ids []uuid.UUID, target uuid.UUID) bool {
	for _, id := range ids {
		if id == target {
			return true
		}
	}
	return false
}

func TestNewAppointment(t *testing.T) {
	t.Run("creates appointment with given fields", func(t *testing.T) {
		userID := uuid.NewV7()
		masterID := uuid.NewV7()
		serviceIDs := []uuid.UUID{uuid.NewV7(), uuid.NewV7()}

		start := time.Date(2026, time.September, 18, 12, 0, 0, 0, time.UTC)
		end := start.Add(90 * time.Minute)

		appt, err := NewAppointment(userID, masterID, serviceIDs, start, end)
		if err != nil {
			t.Fatalf("NewAppointment() unexpected error: %v", err)
		}

		if appt == nil {
			t.Fatal("NewAppointment() returned nil")
		}

		var zeroID uuid.UUID
		if appt.id == zeroID {
			t.Error("expected generated appointment id to be non-zero")
		}

		if appt.userID != userID {
			t.Errorf("userID = %v, want %v", appt.userID, userID)
		}

		if appt.masterID != masterID {
			t.Errorf("masterID = %v, want %v", appt.masterID, masterID)
		}

		if len(appt.serviceIDs) != len(serviceIDs) {
			t.Fatalf("serviceIDs length = %d, want %d", len(appt.serviceIDs), len(serviceIDs))
		}

		for i, id := range serviceIDs {
			if appt.serviceIDs[i] != id {
				t.Errorf("serviceIDs[%d] = %v, want %v", i, appt.serviceIDs[i], id)
			}
		}

		if !appt.timeStart.Equal(start) {
			t.Errorf("timeStart = %v, want %v", appt.timeStart, start)
		}

		if !appt.timeEnd.Equal(end) {
			t.Errorf("timeEnd = %v, want %v", appt.timeEnd, end)
		}

		if appt.createdAt.IsZero() {
			t.Error("createdAt is zero, want non-zero")
		}
	})

	t.Run("creates appointment with empty services", func(t *testing.T) {
		appt := mustNewAppointment(t, []uuid.UUID{})

		if len(appt.serviceIDs) != 0 {
			t.Errorf("serviceIDs length = %d, want 0", len(appt.serviceIDs))
		}
	})

	t.Run("creates appointment with nil services and allows adding later", func(t *testing.T) {
		appt := mustNewAppointment(t, nil)

		if len(appt.serviceIDs) != 0 {
			t.Fatalf("serviceIDs length = %d, want 0", len(appt.serviceIDs))
		}

		serviceID := uuid.NewV7()

		if err := appt.AddIDService(serviceID); err != nil {
			t.Fatalf("AddIDService() unexpected error: %v", err)
		}

		if len(appt.serviceIDs) != 1 {
			t.Fatalf("serviceIDs length = %d, want 1", len(appt.serviceIDs))
		}

		if !containsUUID(appt.serviceIDs, serviceID) {
			t.Errorf("serviceIDs = %v, want contains %v", appt.serviceIDs, serviceID)
		}
	})

	t.Run("generates different ids for different appointments", func(t *testing.T) {
		first := mustNewAppointment(t, nil)
		second := mustNewAppointment(t, nil)

		if first.id == second.id {
			t.Errorf("expected different appointment ids, got %v", first.id)
		}
	})
}

func TestAddIDService(t *testing.T) {
	t.Run("adds new service id", func(t *testing.T) {
		existingID := uuid.NewV7()
		newID := uuid.NewV7()

		appt := mustNewAppointment(t, []uuid.UUID{existingID})

		if err := appt.AddIDService(newID); err != nil {
			t.Fatalf("AddIDService() unexpected error: %v", err)
		}

		if len(appt.serviceIDs) != 2 {
			t.Fatalf("serviceIDs length = %d, want 2", len(appt.serviceIDs))
		}

		if !containsUUID(appt.serviceIDs, existingID) {
			t.Errorf("serviceIDs = %v, want contains existing %v", appt.serviceIDs, existingID)
		}

		if !containsUUID(appt.serviceIDs, newID) {
			t.Errorf("serviceIDs = %v, want contains new %v", appt.serviceIDs, newID)
		}
	})

	t.Run("returns ErrAlreadyExist for duplicate service id", func(t *testing.T) {
		firstID := uuid.NewV7()
		secondID := uuid.NewV7()

		appt := mustNewAppointment(t, []uuid.UUID{firstID, secondID})

		before := make([]uuid.UUID, len(appt.serviceIDs))
		copy(before, appt.serviceIDs)

		err := appt.AddIDService(secondID)
		if !errors.Is(err, ErrAlreadyExist) {
			t.Fatalf("AddIDService() error = %v, want %v", err, ErrAlreadyExist)
		}

		if len(appt.serviceIDs) != len(before) {
			t.Fatalf("serviceIDs length changed: got %d, want %d", len(appt.serviceIDs), len(before))
		}

		for i := range before {
			if appt.serviceIDs[i] != before[i] {
				t.Errorf("serviceIDs[%d] changed: got %v, want %v", i, appt.serviceIDs[i], before[i])
			}
		}
	})

	t.Run("adds service id to nil slice", func(t *testing.T) {
		appt := mustNewAppointment(t, nil)

		serviceID := uuid.NewV7()

		if err := appt.AddIDService(serviceID); err != nil {
			t.Fatalf("AddIDService() unexpected error: %v", err)
		}

		if len(appt.serviceIDs) != 1 {
			t.Fatalf("serviceIDs length = %d, want 1", len(appt.serviceIDs))
		}

		if appt.serviceIDs[0] != serviceID {
			t.Errorf("serviceIDs[0] = %v, want %v", appt.serviceIDs[0], serviceID)
		}
	})

	t.Run("adds multiple different service ids", func(t *testing.T) {
		appt := mustNewAppointment(t, nil)

		serviceIDs := []uuid.UUID{
			uuid.NewV7(),
			uuid.NewV7(),
			uuid.NewV7(),
		}

		for _, serviceID := range serviceIDs {
			if err := appt.AddIDService(serviceID); err != nil {
				t.Fatalf("AddIDService(%v) unexpected error: %v", serviceID, err)
			}
		}

		if len(appt.serviceIDs) != len(serviceIDs) {
			t.Fatalf("serviceIDs length = %d, want %d", len(appt.serviceIDs), len(serviceIDs))
		}

		for _, serviceID := range serviceIDs {
			if !containsUUID(appt.serviceIDs, serviceID) {
				t.Errorf("serviceIDs = %v, want contains %v", appt.serviceIDs, serviceID)
			}
		}
	})

	t.Run("treats zero uuid as regular value", func(t *testing.T) {
		appt := mustNewAppointment(t, nil)

		var zeroID uuid.UUID

		if err := appt.AddIDService(zeroID); err != nil {
			t.Fatalf("AddIDService(zero) unexpected error: %v", err)
		}

		err := appt.AddIDService(zeroID)
		if !errors.Is(err, ErrAlreadyExist) {
			t.Fatalf("AddIDService(zero) second call error = %v, want %v", err, ErrAlreadyExist)
		}
	})
}

func TestTimeStart(t *testing.T) {
	t.Run("returns stored start time", func(t *testing.T) {
		start := time.Date(2026, time.September, 18, 14, 30, 0, 0, time.UTC)
		end := start.Add(time.Hour)

		appt := mustNewAppointmentWithTimes(t, nil, start, end)

		if got := appt.TimeStart(); !got.Equal(start) {
			t.Errorf("TimeStart() = %v, want %v", got, start)
		}
	})

	t.Run("returns zero time for zero appointment value", func(t *testing.T) {
		var appt Appointment

		if got := appt.TimeStart(); !got.IsZero() {
			t.Errorf("TimeStart() = %v, want zero time", got)
		}
	})
}

func TestTimeEnd(t *testing.T) {
	t.Run("returns stored end time", func(t *testing.T) {
		start := time.Date(2026, time.September, 18, 14, 30, 0, 0, time.UTC)
		end := start.Add(2 * time.Hour)

		appt := mustNewAppointmentWithTimes(t, nil, start, end)

		if got := appt.TimeEnd(); !got.Equal(end) {
			t.Errorf("TimeEnd() = %v, want %v", got, end)
		}
	})

	t.Run("returns zero time for zero appointment value", func(t *testing.T) {
		var appt Appointment

		if got := appt.TimeEnd(); !got.IsZero() {
			t.Errorf("TimeEnd() = %v, want zero time", got)
		}
	})
}

func TestAppointmentWorkflow(t *testing.T) {
	userID := uuid.NewV7()
	masterID := uuid.NewV7()

	firstServiceID := uuid.NewV7()
	secondServiceID := uuid.NewV7()
	thirdServiceID := uuid.NewV7()

	start := time.Date(2026, time.September, 20, 11, 0, 0, 0, time.UTC)
	end := start.Add(time.Hour)

	appt, err := NewAppointment(userID, masterID, []uuid.UUID{firstServiceID}, start, end)
	if err != nil {
		t.Fatalf("NewAppointment() unexpected error: %v", err)
	}

	if err := appt.AddIDService(secondServiceID); err != nil {
		t.Fatalf("AddIDService(second) unexpected error: %v", err)
	}

	if err := appt.AddIDService(thirdServiceID); err != nil {
		t.Fatalf("AddIDService(third) unexpected error: %v", err)
	}

	err = appt.AddIDService(secondServiceID)
	if !errors.Is(err, ErrAlreadyExist) {
		t.Fatalf("AddIDService(duplicate) error = %v, want %v", err, ErrAlreadyExist)
	}

	if len(appt.serviceIDs) != 3 {
		t.Fatalf("serviceIDs length = %d, want 3", len(appt.serviceIDs))
	}

	for _, serviceID := range []uuid.UUID{firstServiceID, secondServiceID, thirdServiceID} {
		if !containsUUID(appt.serviceIDs, serviceID) {
			t.Errorf("serviceIDs = %v, want contains %v", appt.serviceIDs, serviceID)
		}
	}

	if !appt.TimeStart().Equal(start) {
		t.Errorf("TimeStart() = %v, want %v", appt.TimeStart(), start)
	}

	if !appt.TimeEnd().Equal(end) {
		t.Errorf("TimeEnd() = %v, want %v", appt.TimeEnd(), end)
	}
}
