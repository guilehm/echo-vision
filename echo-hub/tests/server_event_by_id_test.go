package tests

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/guilehm/echo-vision/echo-hub/internal/app/domain"
	"github.com/guilehm/echo-vision/echo-hub/internal/app/ports/dtos"
	"github.com/guilehm/echo-vision/echo-hub/internal/infra/web"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Event By ID Handler", func() {
	var u *domain.User
	BeforeEach(func() {
		u = saveUser(makeUser("arthur@gmail.com"))
	})

	Context("Get Event By ID", func() {
		It("should return an event by ID successfully", func() {
			// Arrange
			event := saveEvent(makeEvent(u))

			req, err := http.NewRequest(
				http.MethodGet,
				fmt.Sprintf("%s/events/%s", server.URL, event.ID().String()),
				nil,
			)
			Expect(err).ToNot(HaveOccurred())
			token := u.AccessToken()
			req.Header.Set("Authorization", token)

			// Act
			client := http.Client{}
			resp, err := client.Do(req)
			Expect(err).ToNot(HaveOccurred())
			defer resp.Body.Close()

			// Assert
			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			var apiResp web.ApiResponse[dtos.EventResponse]
			err = json.NewDecoder(resp.Body).Decode(&apiResp)
			Expect(err).ToNot(HaveOccurred())
			Expect(apiResp.Data).ToNot(BeNil())
			Expect(apiResp.Data.ID).To(Equal(event.ID()))
			Expect(apiResp.Data.EventType).To(Equal(event.EventType().String()))
			Expect(apiResp.Data.SubType).To(Equal(event.SubType().String()))
			Expect(apiResp.Data.Status).To(Equal(event.Status().String()))
			Expect(apiResp.Data.UserID).To(Equal(event.UserID()))
			Expect(apiResp.Error).To(BeEmpty())
		})

		It("should return an event with file by ID successfully", func() {
			// Arrange
			event := saveEvent(validEventWithFile)

			req, err := http.NewRequest(
				http.MethodGet,
				fmt.Sprintf("%s/events/%s", server.URL, event.ID().String()),
				nil,
			)
			Expect(err).ToNot(HaveOccurred())
			token := u.AccessToken()
			req.Header.Set("Authorization", token)

			// Act
			client := http.Client{}
			resp, err := client.Do(req)
			Expect(err).ToNot(HaveOccurred())
			defer resp.Body.Close()

			// Assert
			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			var apiResp web.ApiResponse[dtos.EventResponse]
			err = json.NewDecoder(resp.Body).Decode(&apiResp)
			Expect(err).ToNot(HaveOccurred())
			Expect(apiResp.Data).ToNot(BeNil())
			Expect(apiResp.Data.ID).To(Equal(event.ID()))
			Expect(apiResp.Data.EventType).To(Equal(event.EventType().String()))
			Expect(apiResp.Data.SubType).To(Equal(event.SubType().String()))
			Expect(apiResp.Data.Status).To(Equal(event.Status().String()))
			Expect(apiResp.Data.UserID).To(Equal(event.UserID()))
			Expect(apiResp.Data.File).ToNot(BeNil())
			Expect(apiResp.Data.File.ContentType).To(Equal(event.File().ContentType))
			Expect(apiResp.Data.File.Filename).To(Equal(event.File().Filename))
			Expect(apiResp.Data.File.ContentType).To(Equal(event.File().ContentType))
			Expect(apiResp.Data.File.Filesize).To(Equal(event.File().Filesize))
			Expect(apiResp.Data.File.URL).To(ContainSubstring(event.File().Filepath))
			Expect(apiResp.Error).To(BeEmpty())
		})

		It("should return 404 if event is not found", func() {
			// Arrange
			notFoundID := uuid.New()
			req, err := http.NewRequest(
				http.MethodGet,
				fmt.Sprintf("%s/events/%s", server.URL, notFoundID.String()),
				nil,
			)
			Expect(err).ToNot(HaveOccurred())
			token := u.AccessToken()
			req.Header.Set("Authorization", token)

			// Act
			client := http.Client{}
			resp, err := client.Do(req)
			Expect(err).ToNot(HaveOccurred())
			defer resp.Body.Close()

			// Assert
			Expect(resp.StatusCode).To(Equal(http.StatusNotFound))

			var apiResp web.ApiResponse[any]
			err = json.NewDecoder(resp.Body).Decode(&apiResp)
			Expect(err).ToNot(HaveOccurred())
			Expect(apiResp.Data).To(BeNil())
			Expect(apiResp.Error).To(Equal("Event not found"))
		})

		It("should return 400 for an invalid event ID format", func() {
			// Arrange
			invalidID := "invalid-uuid-format"
			req, err := http.NewRequest(
				http.MethodGet,
				fmt.Sprintf("%s/events/%s", server.URL, invalidID),
				nil,
			)
			Expect(err).ToNot(HaveOccurred())
			token := u.AccessToken()
			req.Header.Set("Authorization", token)

			// Act
			client := http.Client{}
			resp, err := client.Do(req)
			Expect(err).ToNot(HaveOccurred())
			defer resp.Body.Close()

			// Assert
			Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))

			var apiResp web.ApiResponse[any]
			err = json.NewDecoder(resp.Body).Decode(&apiResp)
			Expect(err).ToNot(HaveOccurred())
			Expect(apiResp.Data).To(BeNil())
			Expect(apiResp.Error).To(Equal("Invalid event ID"))
		})

		It("should return 403 if not authenticated", func() {
			// Arrange
			event := saveEvent(makeEvent(u))

			req, err := http.NewRequest(
				http.MethodGet,
				fmt.Sprintf("%s/events/%s", server.URL, event.ID().String()),
				nil,
			)
			Expect(err).ToNot(HaveOccurred())

			// Act
			client := http.Client{}
			resp, err := client.Do(req)
			Expect(err).ToNot(HaveOccurred())
			defer resp.Body.Close()

			// Assert
			Expect(resp.StatusCode).To(Equal(http.StatusForbidden))

			var apiResp web.ApiResponse[any]
			err = json.NewDecoder(resp.Body).Decode(&apiResp)
			Expect(err).ToNot(HaveOccurred())
			Expect(apiResp.Data).To(BeNil())
			Expect(apiResp.Error).ToNot(BeEmpty())
		})
	})
})
