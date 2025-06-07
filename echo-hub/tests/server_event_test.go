package tests

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/guilehm/echo-vision/echo-hub/internal/app/domain"
	"github.com/guilehm/echo-vision/echo-hub/internal/app/ports"
	"github.com/guilehm/echo-vision/echo-hub/internal/app/ports/dtos"
	"github.com/guilehm/echo-vision/echo-hub/internal/infra/postgres/generated/ent/event"
	"github.com/guilehm/echo-vision/echo-hub/internal/infra/postgres/generated/ent/file"
	"github.com/guilehm/echo-vision/echo-hub/internal/infra/web"
	hubevents "github.com/guilehm/echo-vision/echo-hub/pkg/events"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Event Handler", func() {
	var u *domain.User
	BeforeEach(func() {
		u = saveUser(makeUser("arthur@gmail.com"))
	})

	Context("Create Event", func() {
		It("should create a event successfully", func() {
			// Arrange
			input := ports.EventCreateInput{
				EventType: hubevents.EventTypeImageAnalysis.String(),
				SubType:   hubevents.EventSubTypeDetectLabels.String(),
			}

			req, err := http.NewRequest(
				http.MethodPost,
				fmt.Sprintf("%s/events", server.URL),
				toReader(input),
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

			var apiResp web.ApiResponse[dtos.EventCreateResponse]
			err = json.NewDecoder(resp.Body).Decode(&apiResp)
			Expect(err).ToNot(HaveOccurred())
			Expect(apiResp.Data).ToNot(BeNil())
			Expect(apiResp.Data.ID).ToNot(BeNil())
			Expect(apiResp.Error).To(BeEmpty())

			createdEvent := entClient.Event.Query().
				Where(event.IDEQ(apiResp.Data.ID)).
				OnlyX(ctx)

			Expect(createdEvent.Type).To(BeEquivalentTo(hubevents.EventTypeImageAnalysis))
			Expect(createdEvent.SubType).To(BeEquivalentTo(hubevents.EventSubTypeDetectLabels))
			Expect(createdEvent.Status).To(BeEquivalentTo(hubevents.EventStatusPending))
			Expect(createdEvent.UserID.String()).To(BeEquivalentTo(u.ID().String()))
		})

		It("should return 401 for an invalid token", func() {
			// Arrange
			input := ports.EventCreateInput{
				EventType: hubevents.EventTypeImageAnalysis.String(),
				SubType:   hubevents.EventSubTypeDetectLabels.String(),
			}

			req, err := http.NewRequest(
				http.MethodPost,
				fmt.Sprintf("%s/events", server.URL),
				toReader(input),
			)
			Expect(err).ToNot(HaveOccurred())
			req.Header.Set("Authorization", "InvalidToken")

			// Act
			client := http.Client{}
			resp, err := client.Do(req)
			Expect(err).ToNot(HaveOccurred())
			defer resp.Body.Close()

			// Assert
			Expect(resp.StatusCode).To(Equal(http.StatusUnauthorized))
		})

		It("should return 403 for a missing token", func() {
			// Arrange
			input := ports.EventCreateInput{
				EventType: hubevents.EventTypeImageAnalysis.String(),
				SubType:   hubevents.EventSubTypeDetectLabels.String(),
			}

			req, err := http.NewRequest(
				http.MethodPost,
				fmt.Sprintf("%s/events", server.URL),
				toReader(input),
			)
			Expect(err).ToNot(HaveOccurred())

			// Act
			client := http.Client{}
			resp, err := client.Do(req)
			Expect(err).ToNot(HaveOccurred())
			defer resp.Body.Close()

			// Assert
			Expect(resp.StatusCode).To(Equal(http.StatusForbidden))
		})

		It("should return 400 for an invalid EventType", func() {
			// Arrange
			input := ports.EventCreateInput{
				EventType: "InvalidType",
				SubType:   hubevents.EventSubTypeDetectLabels.String(),
			}

			req, err := http.NewRequest(
				http.MethodPost,
				fmt.Sprintf("%s/events", server.URL),
				toReader(input),
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
			Expect(apiResp.Error).ToNot(BeEmpty())
		})

		It("should return 400 for an invalid SubType", func() {
			// Arrange
			input := ports.EventCreateInput{
				EventType: hubevents.EventTypeImageAnalysis.String(),
				SubType:   "InvalidSubType",
			}

			req, err := http.NewRequest(
				http.MethodPost,
				fmt.Sprintf("%s/events", server.URL),
				toReader(input),
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
			Expect(apiResp.Error).ToNot(BeEmpty())
		})

		Context("Create Event with File", func() {
			It("should create a event with a file successfully", func() {
				// Arrange
				input := ports.EventCreateInput{
					EventType:   hubevents.EventTypeImageAnalysis.String(),
					SubType:     hubevents.EventSubTypeDetectLabels.String(),
					Filepath:    "path/to/file.jpg",
					Filename:    "file.jpg",
					ContentType: "image/jpeg",
					Filesize:    1024,
				}

				req, err := http.NewRequest(
					http.MethodPost,
					fmt.Sprintf("%s/events", server.URL),
					toReader(input),
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

				createdFile := entClient.File.Query().
					Where(
						file.FilenameEQ(input.Filename),
						file.Filepath(input.Filepath),
						file.ContentType(input.ContentType),
						file.FilesizeEQ(input.Filesize),
					).
					OnlyX(ctx)

				Expect(createdFile).ToNot(BeNil())
				Expect(createdFile.Filename).To(BeEquivalentTo(input.Filename))
				Expect(createdFile.Filepath).To(BeEquivalentTo(input.Filepath))
				Expect(createdFile.ContentType).To(BeEquivalentTo(input.ContentType))
				Expect(createdFile.Filesize).To(BeEquivalentTo(input.Filesize))
			})

			It("should return 400 for an invalid file", func() {
				// Arrange
				input := ports.EventCreateInput{
					EventType:   hubevents.EventTypeImageAnalysis.String(),
					SubType:     hubevents.EventSubTypeDetectLabels.String(),
					Filepath:    "path/to/file",
					Filename:    "file.jpg",
					ContentType: "",
					Filesize:    1024,
				}

				req, err := http.NewRequest(
					http.MethodPost,
					fmt.Sprintf("%s/events", server.URL),
					toReader(input),
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
			})
		})
	})

	Context("List Events", func() {
		var u2 *domain.User
		BeforeEach(func() {
			u2 = saveUser(makeUser("dutch@gmail.com"))
			saveEvent(makeEvent(u))
			saveEvent(makeEvent(u))
			saveEvent(makeEvent(u2))
		})

		It("should list events successfully", func() {
			// Arrange
			req, err := http.NewRequest(
				http.MethodGet,
				fmt.Sprintf("%s/users/%s/events", server.URL, u.ID().String()),
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

			var apiResp web.ApiResponse[ports.ApiListResponse[dtos.EventResponse]]
			err = json.NewDecoder(resp.Body).Decode(&apiResp)
			Expect(err).ToNot(HaveOccurred())
			Expect(apiResp.Data).ToNot(BeNil())
			Expect(apiResp.Error).To(BeEmpty())
			Expect(len(apiResp.Data.Results)).To(Equal(2))
			for i := range apiResp.Data.Results {
				Expect(apiResp.Data.Results[i].ID).ToNot(BeNil())
				Expect(apiResp.Data.Results[i].EventType).ToNot(BeEmpty())
				Expect(apiResp.Data.Results[i].SubType).ToNot(BeEmpty())
				Expect(apiResp.Data.Results[i].Status).ToNot(BeEmpty())
				Expect(apiResp.Data.Results[i].UserID).ToNot(BeEmpty())
				Expect(apiResp.Data.Results[i].UserID.String()).To(Equal(u.ID().String()))
			}
		})

		It("should list own events successfully", func() {
			// Arrange
			req, err := http.NewRequest(
				http.MethodGet,
				fmt.Sprintf("%s/events", server.URL),
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

			var apiResp web.ApiResponse[ports.ApiListResponse[dtos.EventResponse]]
			err = json.NewDecoder(resp.Body).Decode(&apiResp)
			Expect(err).ToNot(HaveOccurred())
			Expect(apiResp.Data).ToNot(BeNil())
			Expect(apiResp.Error).To(BeEmpty())
			Expect(len(apiResp.Data.Results)).To(Equal(2))
			for i := range apiResp.Data.Results {
				Expect(apiResp.Data.Results[i].ID).ToNot(BeNil())
				Expect(apiResp.Data.Results[i].EventType).ToNot(BeEmpty())
				Expect(apiResp.Data.Results[i].SubType).ToNot(BeEmpty())
				Expect(apiResp.Data.Results[i].Status).ToNot(BeEmpty())
				Expect(apiResp.Data.Results[i].UserID).ToNot(BeEmpty())
				Expect(apiResp.Data.Results[i].UserID.String()).To(Equal(u.ID().String()))
			}
		})

		It("should not list own events if not authenticated", func() {
			// Arrange
			req, err := http.NewRequest(
				http.MethodGet,
				fmt.Sprintf("%s/events", server.URL),
				nil,
			)
			Expect(err).ToNot(HaveOccurred())

			// Do not set the Authorization header

			// Act
			client := http.Client{}
			resp, err := client.Do(req)
			Expect(err).ToNot(HaveOccurred())
			defer resp.Body.Close()

			// Assert
			Expect(resp.StatusCode).To(Equal(http.StatusForbidden))

			var apiResp web.ApiResponse[ports.ApiListResponse[dtos.EventResponse]]
			err = json.NewDecoder(resp.Body).Decode(&apiResp)
			Expect(err).ToNot(HaveOccurred())
			Expect(apiResp.Data).To(BeNil())
			Expect(apiResp.Error).ToNot(BeEmpty())
		})
		It("should not list events from another user", func() {
			// Arrange
			req, err := http.NewRequest(
				http.MethodGet,
				fmt.Sprintf("%s/users/%s/events", server.URL, u.ID().String()),
				nil,
			)
			Expect(err).ToNot(HaveOccurred())
			// Use u2's token
			token := u2.AccessToken()
			req.Header.Set("Authorization", token)

			// Act
			client := http.Client{}
			resp, err := client.Do(req)
			Expect(err).ToNot(HaveOccurred())
			defer resp.Body.Close()

			// Assert
			Expect(resp.StatusCode).To(Equal(http.StatusUnauthorized))

			var apiResp web.ApiResponse[ports.ApiListResponse[dtos.EventResponse]]
			err = json.NewDecoder(resp.Body).Decode(&apiResp)
			Expect(err).ToNot(HaveOccurred())
			Expect(apiResp.Data).To(BeNil())
			Expect(apiResp.Error).ToNot(BeEmpty())
		})
	})
})
