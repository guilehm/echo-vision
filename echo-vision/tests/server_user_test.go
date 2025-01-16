package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/guilehm/echo-vision/internal/app/ports"
	"github.com/guilehm/echo-vision/internal/app/shared"
	"github.com/guilehm/echo-vision/internal/infra/postgres/generated/ent/user"
	"github.com/guilehm/echo-vision/internal/infra/web"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("User Handler", func() {
	Context("CreateUser", func() {
		It("should create a user successfully", func() {
			// Arrange
			input := ports.UserCreateInput{
				FirstName: validUser.FirstName(),
				LastName:  validUser.LastName(),
				Email:     strings.ToUpper(validUser.Email()),
				Password:  "awesomepassword",
			}
			b, _ := json.Marshal(input)

			// Act
			resp, err := http.Post(fmt.Sprintf("%s/users", server.URL), "application/json", bytes.NewReader(b))
			Expect(err).ToNot(HaveOccurred())
			defer resp.Body.Close()

			// Assert
			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			var apiResp web.ApiResponse[ports.UserCreateResponse]
			err = json.NewDecoder(resp.Body).Decode(&apiResp)
			Expect(err).ToNot(HaveOccurred())
			Expect(apiResp.Data).ToNot(BeNil())
			Expect(apiResp.Data.ID).ToNot(BeNil())
			Expect(apiResp.Error).To(BeEmpty())

			createdUser := entClient.User.Query().
				Where(user.IDEQ(apiResp.Data.ID)).
				OnlyX(ctx)

			Expect(createdUser.FirstName).To(Equal(validUser.FirstName()))
			Expect(createdUser.LastName).To(Equal(validUser.LastName()))
			Expect(createdUser.Email).To(Equal(validUser.Email()))
			Expect(createdUser.Password).ToNot(BeEmpty())
			Expect(createdUser.Password).ToNot(Equal(input.Password))
			Expect(createdUser.AccessToken).ToNot(BeEmpty())
			Expect(createdUser.RefreshToken).ToNot(BeEmpty())

			passwordErr := passwordAdapter.ValidatePassword(input.Password, createdUser.Password)
			Expect(passwordErr).ToNot(HaveOccurred())
		})

		It("should return 400 error for invalid password", func() {
			// Arrange
			input := ports.UserCreateInput{
				FirstName: validUser.FirstName(),
				LastName:  validUser.LastName(),
				Email:     strings.ToUpper(validUser.Email()),
				Password:  "pass",
			}
			b, _ := json.Marshal(input)

			// Act
			resp, err := http.Post(fmt.Sprintf("%s/users", server.URL), "application/json", bytes.NewReader(b))
			Expect(err).ToNot(HaveOccurred())
			defer resp.Body.Close()

			// Assert
			Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))

			var apiResp web.ApiResponse[ports.UserCreateResponse]
			err = json.NewDecoder(resp.Body).Decode(&apiResp)
			Expect(err).ToNot(HaveOccurred())
			Expect(apiResp.Data).To(BeNil())
			Expect(apiResp.Error).ToNot(BeEmpty())
			Expect(apiResp.Error).To(Equal(shared.ErrInvalidPassword.Error()))

			count := entClient.User.Query().
				CountX(ctx)

			Expect(count).To(Equal(0))
		})

		It("should return 400 error if email already exists", func() {
			// Arrange
			input := ports.UserCreateInput{
				FirstName: validUser.FirstName(),
				LastName:  validUser.LastName(),
				Email:     validUser.Email(),
			}
			b, _ := json.Marshal(input)
			_, err := repo.SaveUser(ctx, nil, validUser)
			Expect(err).ToNot(HaveOccurred())

			// Act
			resp, err := http.Post(fmt.Sprintf("%s/users", server.URL), "application/json", bytes.NewReader(b))
			Expect(err).ToNot(HaveOccurred())
			defer resp.Body.Close()

			// Assert
			Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))

			var apiResp web.ApiResponse[any]
			err = json.NewDecoder(resp.Body).Decode(&apiResp)
			Expect(err).ToNot(HaveOccurred())
			Expect(apiResp.Data).To(BeNil())
			Expect(apiResp.Error).ToNot(BeEmpty())
		})
		It("should return 400 error if email is invalid", func() {
			// Arrange
			input := ports.UserCreateInput{
				FirstName: validUser.FirstName(),
				LastName:  validUser.LastName(),
				Email:     "invalid-email",
			}
			b, _ := json.Marshal(input)

			// Act
			resp, err := http.Post(fmt.Sprintf("%s/users", server.URL), "application/json", bytes.NewReader(b))
			Expect(err).ToNot(HaveOccurred())
			defer resp.Body.Close()

			// Assert
			Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))

			var apiResp web.ApiResponse[any]
			err = json.NewDecoder(resp.Body).Decode(&apiResp)
			Expect(err).ToNot(HaveOccurred())
			Expect(apiResp.Data).To(BeNil())
			Expect(apiResp.Error).To(Equal(shared.ErrInvalidEmail.Error()))
		})
		It("should return 400 error if payload is invalid", func() {
			// Arrange
			input := struct {
				Error string `json:"error"`
			}{
				Error: "invalid payload",
			}
			b, _ := json.Marshal(input)

			// Act
			resp, err := http.Post(fmt.Sprintf("%s/users", server.URL), "application/json", bytes.NewReader(b))
			Expect(err).ToNot(HaveOccurred())
			defer resp.Body.Close()

			// Assert
			Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))

			var apiResp web.ApiResponse[any]
			err = json.NewDecoder(resp.Body).Decode(&apiResp)
			Expect(err).ToNot(HaveOccurred())
			Expect(apiResp.Data).To(BeNil())
			Expect(apiResp.Error).ToNot(BeEmpty())
		})
	})
})
