package tests

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/guilehm/echo-vision/echo-hub/internal/app/ports"
	"github.com/guilehm/echo-vision/echo-hub/internal/infra/web"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Uploads Handler", func() {
	Context("when creating a presigned URL", func() {
		It("should return a presigned URL", func() {
			// Arrange

			// create user
			u := saveUser(makeUser("mario@nintendo.com"))

			// create input
			input := ports.UploadPresignedURLInput{
				Filename:    "test.jpg",
				Filepath:    "users/1234/abcd.jpg",
				ContentType: "image/jpeg",
			}
			req, err := http.NewRequest(
				http.MethodPost,
				fmt.Sprintf("%s/uploads/presigned-url", server.URL),
				toReader(input),
			)
			Expect(err).ToNot(HaveOccurred())
			// authenticate user
			req.Header.Set("Authorization", u.AccessToken())

			// Act
			client := http.Client{}
			resp, err := client.Do(req)
			Expect(err).ToNot(HaveOccurred())
			defer resp.Body.Close()

			// Assert
			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			var apiResp web.ApiResponse[ports.UploadPresignedURLResponse]
			err = json.NewDecoder(resp.Body).Decode(&apiResp)
			Expect(err).ToNot(HaveOccurred())

			Expect(apiResp.Data.URL).ToNot(BeEmpty())
		})

		It("should return 403 for not authenticated user", func() {
			// Arrange
			input := ports.UploadPresignedURLInput{
				Filename:    "test.jpg",
				Filepath:    "users/1234/abcd.jpg",
				ContentType: "image/jpeg",
			}

			// Act
			resp, err := http.Post(
				fmt.Sprintf("%s/uploads/presigned-url", server.URL),
				"application/json",
				toReader(input),
			)
			Expect(err).ToNot(HaveOccurred())
			defer resp.Body.Close()

			// Assert
			Expect(resp.StatusCode).To(Equal(http.StatusForbidden))

			var apiResp web.ApiResponse[ports.UploadPresignedURLResponse]
			err = json.NewDecoder(resp.Body).Decode(&apiResp)
			Expect(err).ToNot(HaveOccurred())

			Expect(apiResp.Data).To(BeNil())
			Expect(apiResp.Error).ToNot(BeNil())
		})

		It("should return 400 for invalid input", func() {
			// Arrange
			u := saveUser(makeUser("mario@nintendo.com"))
			input := ports.UploadPresignedURLInput{
				Filename:    "",
				Filepath:    "",
				ContentType: "",
			}

			// Act
			req, err := http.NewRequest(
				http.MethodPost,
				fmt.Sprintf("%s/uploads/presigned-url", server.URL),
				toReader(input),
			)
			Expect(err).ToNot(HaveOccurred())

			// authenticate user
			req.Header.Set("Authorization", u.AccessToken())

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
