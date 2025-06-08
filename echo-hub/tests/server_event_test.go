package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/guilehm/echo-vision/echo-hub/internal/app/domain"
	"github.com/guilehm/echo-vision/echo-hub/internal/app/ports"
	"github.com/guilehm/echo-vision/echo-hub/internal/app/ports/dtos"
	"github.com/guilehm/echo-vision/echo-hub/internal/infra/postgres"
	"github.com/guilehm/echo-vision/echo-hub/internal/infra/postgres/generated/ent"
	"github.com/guilehm/echo-vision/echo-hub/internal/infra/postgres/generated/ent/event"
	"github.com/guilehm/echo-vision/echo-hub/internal/infra/postgres/generated/ent/file"
	"github.com/guilehm/echo-vision/echo-hub/internal/infra/web"
	hubevents "github.com/guilehm/echo-vision/echo-hub/pkg/events"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Event Handler", func() {
	var u *domain.User
	topic := hubevents.EventImageAnalysCreated
	BeforeEach(func() {
		u = saveUser(makeUser("arthur@gmail.com"))
	})

	Context("Create Event", func() {
		It("should create a event successfully", func() {
			// Arrange

			done := make(chan any)
			expectFunc := func(msg hubevents.EventMessage) {
				Expect(msg).ToNot(BeNil())
				Expect(msg.ID).ToNot(BeNil())
				Expect(msg.Type).To(Equal(hubevents.EventTypeImageAnalysis))
				Expect(msg.SubType).To(Equal(hubevents.EventSubTypeDetectLabels))
				Expect(msg.File).To(BeNil())
			}
			expectMessageCalled(topic, done, expectFunc)

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
			Eventually(done, "1s").Should(BeClosed())
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

				done := make(chan any)
				expectFunc := func(msg hubevents.EventMessage) {
					Expect(msg).ToNot(BeNil())
					Expect(msg.ID).ToNot(BeNil())
					Expect(msg.Type).To(Equal(hubevents.EventTypeImageAnalysis))
					Expect(msg.SubType).To(Equal(hubevents.EventSubTypeDetectLabels))
					Expect(msg.File).ToNot(BeNil())
				}
				expectMessageCalled(topic, done, expectFunc)

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
				Eventually(done, "1s").Should(BeClosed())
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

		FIt("should list paginated events", func() {
			// Arrange
			firstEvent := entClient.Event.Query().
				Where(event.UserIDEQ(u.ID())).
				Order(ent.Desc(event.FieldCreatedAt)).
				FirstX(ctx)

			// add pagination param
			paginationLimit := 1
			params := url.Values{}
			params.Add("limit", fmt.Sprintf("%d", paginationLimit))

			req, err := http.NewRequest(
				http.MethodGet,
				fmt.Sprintf("%s/users/%s/events", server.URL, u.ID().String()),
				nil,
			)
			Expect(err).ToNot(HaveOccurred())
			token := u.AccessToken()
			req.Header.Set("Authorization", token)
			req.URL.RawQuery = params.Encode()

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
			Expect(len(apiResp.Data.Results)).To(Equal(1))
			for i := range apiResp.Data.Results {
				Expect(apiResp.Data.Results[i].ID).To(Equal(firstEvent.ID))
			}
			Expect(apiResp.Data.NextCursor).ToNot(BeEmpty())
			Expect(apiResp.Data.NextCursor).To(Equal(postgres.EncodeCursorForTest(firstEvent.CreatedAt, firstEvent.ID)))
		})

		It("should list paginated events with cursor", func() {
			// Arrange
			firstEvent := entClient.Event.Query().
				Where(event.UserIDEQ(u.ID())).
				Order(ent.Desc(event.FieldCreatedAt)).
				FirstX(ctx)

			lastEvent := entClient.Event.Query().
				Where(event.UserIDEQ(u.ID())).
				Order(ent.Asc(event.FieldCreatedAt)).
				FirstX(ctx)

			cursor := postgres.EncodeCursorForTest(firstEvent.CreatedAt, firstEvent.ID)

			// add pagination param
			paginationLimit := 1
			params := url.Values{}
			params.Add("limit", fmt.Sprintf("%d", paginationLimit))
			params.Add("cursor", cursor)

			req, err := http.NewRequest(
				http.MethodGet,
				fmt.Sprintf("%s/users/%s/events", server.URL, u.ID().String()),
				nil,
			)
			Expect(err).ToNot(HaveOccurred())
			token := u.AccessToken()
			req.Header.Set("Authorization", token)
			req.URL.RawQuery = params.Encode()

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
			Expect(len(apiResp.Data.Results)).To(Equal(1))
			for i := range apiResp.Data.Results {
				Expect(apiResp.Data.Results[i].ID).To(Equal(lastEvent.ID))
			}
			Expect(apiResp.Data.NextCursor).To(BeEmpty())
		})

		It("should handle last pagination successfully", func() {
			// Arrange
			lastEvent := entClient.Event.Query().
				Where(event.UserIDEQ(u.ID())).
				Order(ent.Asc(event.FieldCreatedAt)).
				FirstX(ctx)

			cursor := postgres.EncodeCursorForTest(lastEvent.CreatedAt, lastEvent.ID)

			// add pagination param
			paginationLimit := 1
			params := url.Values{}
			params.Add("limit", fmt.Sprintf("%d", paginationLimit))
			params.Add("cursor", cursor)

			req, err := http.NewRequest(
				http.MethodGet,
				fmt.Sprintf("%s/users/%s/events", server.URL, u.ID().String()),
				nil,
			)
			Expect(err).ToNot(HaveOccurred())
			token := u.AccessToken()
			req.Header.Set("Authorization", token)
			req.URL.RawQuery = params.Encode()

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
			Expect(len(apiResp.Data.Results)).To(Equal(1))
			for i := range apiResp.Data.Results {
				Expect(apiResp.Data.Results[i].ID).To(Equal(lastEvent.ID))
			}
			Expect(apiResp.Data.NextCursor).To(BeEmpty())
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

		It("should return 400 for invalid pagination parameters", func() {
			// Arrange
			params := url.Values{}
			params.Add("limit", "invalid")  // Invalid limit
			params.Add("cursor", "invalid") // Invalid cursor
			req, err := http.NewRequest(
				http.MethodGet,
				fmt.Sprintf("%s/users/%s/events?%s", server.URL, u.ID().String(), params.Encode()),
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
			var apiResp web.ApiResponse[ports.ApiListResponse[dtos.EventResponse]]
			err = json.NewDecoder(resp.Body).Decode(&apiResp)
			Expect(err).ToNot(HaveOccurred())
			Expect(apiResp.Data).To(BeNil())
			Expect(apiResp.Error).ToNot(BeEmpty())
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
