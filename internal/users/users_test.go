package users

import (
	"errors"
	"net/mail"
	"reflect"
	"testing"
)

func TestAddUser(t *testing.T) {
	testManager := NewManager()

	testFirstName := "Test"
	testLastName := "Userman"
	testEmail, err := mail.ParseAddress("foo@bar.com")
	if err != nil {
		t.Fatalf("error parsing test email address: %v", err)
	}

	err = testManager.AddUser(testFirstName, testLastName, testEmail.String())
	if err != nil {
		t.Fatalf("error creating user: %v", err)
	}

	if len(testManager.users) != 1 {
		t.Errorf("bad test manager user count, wanted: %d, got %d", 1, len(testManager.users))
		if len(testManager.users) < 1 {
			t.Fatal()
		}
	}

	expectedUser := User{
		FirtName: testFirstName,
		LastName: testLastName,
		Email:    *testEmail,
	}

	foundUser := testManager.users[0]

	if !reflect.DeepEqual(expectedUser, foundUser) {
		t.Errorf("added user data is not correct\nwanted: %+v\ngot: %+v\n", expectedUser, foundUser)
	}
}

func TestAddUserInvalidEmail(t *testing.T) {
	testManager := NewManager()

	testFirstName := "Test"
	testLastName := "Userman"
	testEmail := "foobar"

	err := testManager.AddUser(testFirstName, testLastName, testEmail) //add test values
	if err == nil {                                                    //checking that Adduser func completed without mistake
		t.Error("no error returned for invalid email")
	} else {
		expectedErr := "invalid email: foobar" // checking error
		if err.Error() != expectedErr {
			t.Errorf("bad error text, wanted: %s, got: %s", expectedErr, err.Error())
		}

	}

	if len(testManager.users) > 0 { //checking users count
		t.Errorf("bad test manager user count, wanted: %d, got: %d", 0, len(testManager.users))
	}
}

func TestAddUserEmptyFirstName(t *testing.T) {
	testManager := NewManager()

	testFirstName := "" // Empty firstName
	testLastName := "Userman"
	testEmail, err := mail.ParseAddress("foo@bar.com")
	if err != nil {
		t.Fatalf("error parsing test email address: %v", err)
	}

	err = testManager.AddUser(testFirstName, testLastName, testEmail.String())
	if err == nil {
		t.Error("no error returned for invalid firstname")
	} else {
		expectedErr := "invalid firstname: " //expect empty firstname
		if err.Error() != expectedErr {
			t.Errorf("bad error text, wanted: %s, got: %s", expectedErr, err.Error())
		}
	}

	if len(testManager.users) > 0 { //checking users count
		t.Errorf("bad test manager user count, wanted: %d, got: %d", 0, len(testManager.users))
	}
}

func TestAddUserEmptyLastName(t *testing.T) {
	testManager := NewManager()

	testFirstName := "Test"
	testLastName := "" //empty lastname
	testEmail, err := mail.ParseAddress("foo@bar.com")
	if err != nil {
		t.Fatalf("error parsing test email address: %v", err)
	}

	err = testManager.AddUser(testFirstName, testLastName, testEmail.String())
	if err == nil {
		t.Error("no error returned for invalid lastname")
	} else {
		expectedErr := "invalid lastname: " // expect empty lastname
		if err.Error() != expectedErr {
			t.Errorf("bad error text, wanted: %s, got: %s", expectedErr, err.Error())
		}
	}

	if len(testManager.users) > 0 { //checking users count
		t.Errorf("bad test manager user count, wanted: %d, got: %d", 0, len(testManager.users))
	}
}

func TestAddUserDuplicateName(t *testing.T) {
	testManager := NewManager()

	testFirstName := "Test"
	testLastName := "Userman"
	testEmail, err := mail.ParseAddress("foo@bar.com")
	if err != nil {
		t.Fatalf("error parsing test email address: %v", err)
	}

	err = testManager.AddUser(testFirstName, testLastName, testEmail.String())
	if err != nil {
		t.Fatalf("error parsing test email address: %v", err)
	}

	err = testManager.AddUser(testFirstName, testLastName, testEmail.String())
	if err == nil {
		t.Error("no error returned for duplicate user")
	} else {
		expectedErr := "user with this name already exist"
		if err.Error() != expectedErr {
			t.Errorf("bad error text, wanted: %s, got: %s", expectedErr, err.Error())
		}
	}

	if len(testManager.users) != 1 { //checking users count
		t.Errorf("bad test manager user count, wanted: %d, got: %d", 1, len(testManager.users))
	}
}

func TestGetUserByName(t *testing.T) {
	testManager := NewManager()

	err := testManager.AddUser("foo", "bar", "f.bar@gmail.com")
	if err != nil {
		t.Fatalf("error adding test user: %v", err)
	}

	err = testManager.AddUser("bar", "baz", "barbaz.bar@gmail.com")
	if err != nil {
		t.Fatalf("error adding test user: %v", err)
	}

	err = testManager.AddUser("foo", "baz", "foobaz.bar@gmail.com")
	if err != nil {
		t.Fatalf("error adding test user: %v", err)
	}

	err = testManager.AddUser("baz", "foo", "bazfoo.bar@gmail.com")
	if err != nil {
		t.Fatalf("error adding test user: %v", err)
	}

	tests := map[string]struct {
		first       string
		last        string
		expected    *User
		expectedErr error
	}{
		"simple lookup": {
			first:       "foo",
			last:        "bar",
			expected:    &testManager.users[0],
			expectedErr: nil,
		},
		"last element lookup": {
			first:       "baz",
			last:        "foo",
			expected:    &testManager.users[3],
			expectedErr: nil,
		},
		"no match lookup": {
			first:       "qproop",
			last:        "qlkfmkf",
			expected:    nil,
			expectedErr: ErrNoResultFound,
		},
		"partial match lookup": {
			first:       "foo",
			last:        "foo",
			expected:    nil,
			expectedErr: ErrNoResultFound,
		},
		"empty first name": {
			first:       "",
			last:        "baz",
			expected:    nil,
			expectedErr: ErrNoResultFound,
		},
		"empty last name": {
			first:       "foo",
			last:        "",
			expected:    nil,
			expectedErr: ErrNoResultFound,
		},
	}

	for name, test := range tests {
		result, err := testManager.GetUserByName(test.first, test.last)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("%s: invalid result\ngot: %+v\nwanted: %+v\n", name, result, test.expected)
		}
		if !errors.Is(err, test.expectedErr) {
			t.Errorf("%s: invalid result\ngot: %+v\nwanted: %+v\n", name, result, test.expected)
		}
	}
}
