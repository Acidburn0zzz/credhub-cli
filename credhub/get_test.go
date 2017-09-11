package credhub_test

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/cloudfoundry-incubator/credhub-cli/credhub"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials/values"
	. "github.com/onsi/ginkgo/extensions/table"
)

var _ = Describe("Get", func() {

	Describe("GetLatestVersion()", func() {
		It("requests the credential by name", func() {
			dummyAuth := &DummyAuth{Response: &http.Response{
				Body: ioutil.NopCloser(bytes.NewBufferString("")),
			}}

			ch, _ := New("https://example.com", Auth(dummyAuth.Builder()))

			ch.GetLatestVersion("/example-password")
			url := dummyAuth.Request.URL.String()
			Expect(url).To(Equal("https://example.com/api/v1/data?name=%2Fexample-password&versions=1"))
			Expect(dummyAuth.Request.Method).To(Equal(http.MethodGet))
		})

		Context("when successful", func() {
			It("returns a credential by name", func() {
				responseString := `{
	"data": [
	{
      "id": "some-id",
      "name": "/example-password",
      "type": "password",
      "value": "some-password",
      "version_created_at": "2017-01-05T01:01:01Z"
    }
    ]}`
				dummyAuth := &DummyAuth{Response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(responseString)),
				}}

				ch, _ := New("https://example.com", Auth(dummyAuth.Builder()))

				cred, err := ch.GetLatestVersion("/example-password")
				Expect(err).To(BeNil())
				Expect(cred.Id).To(Equal("some-id"))
				Expect(cred.Name).To(Equal("/example-password"))
				Expect(cred.Type).To(Equal("password"))
				Expect(cred.Value.(string)).To(Equal("some-password"))
				Expect(cred.VersionCreatedAt).To(Equal("2017-01-05T01:01:01Z"))
			})
		})

		Context("when response body cannot be unmarshalled", func() {
			It("returns an error", func() {
				dummyAuth := &DummyAuth{Response: &http.Response{
					Body: ioutil.NopCloser(bytes.NewBufferString("something-invalid")),
				}}

				ch, _ := New("https://example.com", Auth(dummyAuth.Builder()))
				_, err := ch.GetLatestVersion("/example-password")

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when the response body contains an empty list", func() {
			It("returns an error", func() {
				dummyAuth := &DummyAuth{Response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(`{"data":[]}`)),
				}}
				ch, _ := New("https://example.com", Auth(dummyAuth.Builder()))
				_, err := ch.GetLatestVersion("/example-password")

				Expect(err).To(MatchError("response did not contain any credentials"))
			})
		})
	})

	Describe("GetNVersions()", func() {
		It("makes a request for N versions of a credential", func() {
			dummyAuth := &DummyAuth{Response: &http.Response{
				Body: ioutil.NopCloser(bytes.NewBufferString("")),
			}}

			ch, _ := New("https://example.com", Auth(dummyAuth.Builder()))

			ch.GetNVersions("/example-password", 3)
			url := dummyAuth.Request.URL.String()
			Expect(url).To(Equal("https://example.com/api/v1/data?name=%2Fexample-password&versions=3"))
			Expect(dummyAuth.Request.Method).To(Equal(http.MethodGet))
		})

		Context("when successful", func() {
			It("returns a list of N passwords", func() {
				responseString := `{
	"data": [
	{
      "id": "some-id",
      "name": "/example-password",
      "type": "password",
      "value": "some-password",
      "version_created_at": "2017-01-05T01:01:01Z"
    },
	{
      "id": "some-id",
      "name": "/example-password",
      "type": "password",
      "value": "some-other-password",
      "version_created_at": "2017-01-05T01:01:01Z"
    }
    ]}`
				dummyAuth := &DummyAuth{Response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(responseString)),
				}}

				ch, _ := New("https://example.com", Auth(dummyAuth.Builder()))

				creds, err := ch.GetNVersions("/example-password", 2)
				Expect(err).To(BeNil())
				Expect(creds[0].Id).To(Equal("some-id"))
				Expect(creds[0].Name).To(Equal("/example-password"))
				Expect(creds[0].Type).To(Equal("password"))
				Expect(creds[0].Value.(string)).To(Equal("some-password"))
				Expect(creds[0].VersionCreatedAt).To(Equal("2017-01-05T01:01:01Z"))

				Expect(creds[1].Value.(string)).To(Equal("some-other-password"))
			})

			It("returns a list of N users", func() {
				responseString := `{
	"data": [
	{
      "id": "some-id",
      "name": "/example-user",
      "type": "user",
      "value": {
      	"username": "first-username",
      	"password": "dummy_password",
      	"password_hash": "$6$kjhlkjh$lkjhasdflkjhasdflkjh"
      },
      "version_created_at": "2017-01-05T01:01:01Z"
    },
	{
      "id": "some-id",
      "name": "/example-user",
      "type": "user",
      "value": {
      	"username": "second-username",
      	"password": "another_random_dummy_password",
      	"password_hash": "$6$kjhlkjh$lkjhasdflkjhasdflkjh"
      },
      "version_created_at": "2017-01-05T01:01:01Z"
    }
    ]}`
				dummyAuth := &DummyAuth{Response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(responseString)),
				}}

				ch, _ := New("https://example.com", Auth(dummyAuth.Builder()))

				creds, err := ch.GetNVersions("/example-user", 2)
				Expect(err).To(BeNil())
				Expect(creds[0].Id).To(Equal("some-id"))
				Expect(creds[0].Name).To(Equal("/example-user"))
				Expect(creds[0].Type).To(Equal("user"))
				firstCredValue := creds[0].Value.(map[string]interface{})
				Expect(firstCredValue["username"]).To(Equal("first-username"))
				Expect(firstCredValue["password"]).To(Equal("dummy_password"))
				Expect(firstCredValue["password_hash"]).To(Equal("$6$kjhlkjh$lkjhasdflkjhasdflkjh"))
				Expect(creds[0].VersionCreatedAt).To(Equal("2017-01-05T01:01:01Z"))

				Expect(creds[1].Value.(map[string]interface{})["username"]).To(Equal("second-username"))
			})
		})

		Context("when response body cannot be unmarshalled", func() {
			It("returns an error", func() {
				dummyAuth := &DummyAuth{Response: &http.Response{
					Body: ioutil.NopCloser(bytes.NewBufferString("something-invalid")),
				}}

				ch, _ := New("https://example.com", Auth(dummyAuth.Builder()))
				_, err := ch.GetLatestVersion("/example-password")

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when the response body contains an empty list", func() {
			It("returns an error", func() {
				dummyAuth := &DummyAuth{Response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(`{"data":[]}`)),
				}}
				ch, _ := New("https://example.com", Auth(dummyAuth.Builder()))
				_, err := ch.GetLatestVersion("/example-password")

				Expect(err).To(MatchError("response did not contain any credentials"))
			})
		})
	})

	Describe("GetPassword()", func() {
		It("requests the credential by name", func() {
			dummyAuth := &DummyAuth{Response: &http.Response{
				Body: ioutil.NopCloser(bytes.NewBufferString("")),
			}}

			ch, _ := New("https://example.com", Auth(dummyAuth.Builder()))
			ch.GetPassword("/example-password")
			url := dummyAuth.Request.URL
			Expect(url.String()).To(ContainSubstring("https://example.com/api/v1/data"))
			Expect(url.String()).To(ContainSubstring("name=%2Fexample-password"))
			Expect(url.String()).To(ContainSubstring("versions=1"))
			Expect(dummyAuth.Request.Method).To(Equal(http.MethodGet))
		})

		Context("when successful", func() {
			It("returns a password credential", func() {
				responseString := `{
  "data": [
    {
      "id": "some-id",
      "name": "/example-password",
      "type": "password",
      "value": "some-password",
      "version_created_at": "2017-01-05T01:01:01Z"
    }]}`
				dummyAuth := &DummyAuth{Response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(responseString)),
				}}

				ch, _ := New("https://example.com", Auth(dummyAuth.Builder()))
				cred, err := ch.GetPassword("/example-password")
				Expect(err).ToNot(HaveOccurred())
				Expect(cred.Value).To(BeEquivalentTo("some-password"))
			})
		})

		Context("when response body cannot be unmarshalled", func() {
			It("returns an error", func() {
				dummyAuth := &DummyAuth{Response: &http.Response{
					Body: ioutil.NopCloser(bytes.NewBufferString("something-invalid")),
				}}
				ch, _ := New("https://example.com", Auth(dummyAuth.Builder()))
				_, err := ch.GetPassword("/example-cred")

				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("GetCertificate()", func() {
		It("requests the credential by name", func() {
			dummyAuth := &DummyAuth{Response: &http.Response{
				Body: ioutil.NopCloser(bytes.NewBufferString("")),
			}}

			ch, _ := New("https://example.com", Auth(dummyAuth.Builder()))
			ch.GetCertificate("/example-certificate")
			url := dummyAuth.Request.URL
			Expect(url.String()).To(ContainSubstring("https://example.com/api/v1/data"))
			Expect(url.String()).To(ContainSubstring("name=%2Fexample-certificate"))
			Expect(url.String()).To(ContainSubstring("versions=1"))

			Expect(dummyAuth.Request.Method).To(Equal(http.MethodGet))
		})

		Context("when successful", func() {
			It("returns a certificate credential", func() {
				responseString := `{
				  "data": [{
	"id": "some-id",
	"name": "/example-certificate",
	"type": "certificate",
	"value": {
		"ca": "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
		"certificate": "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
		"private_key": "-----BEGIN RSA PRIVATE KEY-----\n...\n-----END RSA PRIVATE KEY-----"
	},
	"version_created_at": "2017-01-01T04:07:18Z"
}]}`
				dummyAuth := &DummyAuth{Response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(responseString)),
				}}

				ch, _ := New("https://example.com", Auth(dummyAuth.Builder()))

				cred, err := ch.GetCertificate("/example-certificate")
				Expect(err).ToNot(HaveOccurred())
				Expect(cred.Value.Ca).To(Equal("-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----"))
				Expect(cred.Value.Certificate).To(Equal("-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----"))
				Expect(cred.Value.PrivateKey).To(Equal("-----BEGIN RSA PRIVATE KEY-----\n...\n-----END RSA PRIVATE KEY-----"))
				Expect(cred.VersionCreatedAt).To(Equal("2017-01-01T04:07:18Z"))
			})
		})

		Context("when response body cannot be unmarshalled", func() {
			It("returns an error", func() {
				dummyAuth := &DummyAuth{Response: &http.Response{
					Body: ioutil.NopCloser(bytes.NewBufferString("something-invalid")),
				}}
				ch, _ := New("https://example.com", Auth(dummyAuth.Builder()))
				_, err := ch.GetCertificate("/example-cred")

				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("GetUser()", func() {
		It("requests the credential by name", func() {
			dummyAuth := &DummyAuth{Response: &http.Response{
				Body: ioutil.NopCloser(bytes.NewBufferString("")),
			}}

			ch, _ := New("https://example.com", Auth(dummyAuth.Builder()))
			ch.GetUser("/example-user")
			url := dummyAuth.Request.URL
			Expect(url.String()).To(ContainSubstring("https://example.com/api/v1/data"))
			Expect(url.String()).To(ContainSubstring("name=%2Fexample-user"))
			Expect(url.String()).To(ContainSubstring("versions=1"))

			Expect(dummyAuth.Request.Method).To(Equal(http.MethodGet))
		})

		Context("when successful", func() {
			It("returns a user credential", func() {
				responseString := `{
				  "data": [
					{
					  "id": "some-id",
					  "name": "/example-user",
					  "type": "user",
					  "value": {
						"username": "some-username",
						"password": "some-password",
						"password_hash": "some-hash"
					  },
					  "version_created_at": "2017-01-05T01:01:01Z"
					}
				  ]
				}`
				dummyAuth := &DummyAuth{Response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(responseString)),
				}}

				ch, _ := New("https://example.com", Auth(dummyAuth.Builder()))
				cred, err := ch.GetUser("/example-user")
				Expect(err).ToNot(HaveOccurred())
				Expect(cred.Value.PasswordHash).To(Equal("some-hash"))
				Expect(cred.Value.User).To(Equal(values.User{
					Username: "some-username",
					Password: "some-password",
				}))
			})
		})

		Context("when response body cannot be unmarshalled", func() {
			It("returns an error", func() {
				dummyAuth := &DummyAuth{Response: &http.Response{
					Body: ioutil.NopCloser(bytes.NewBufferString("something-invalid")),
				}}
				ch, _ := New("https://example.com", Auth(dummyAuth.Builder()))
				_, err := ch.GetUser("/example-cred")

				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("GetRSA()", func() {
		It("requests the credential by name", func() {
			dummyAuth := &DummyAuth{Response: &http.Response{
				Body: ioutil.NopCloser(bytes.NewBufferString("")),
			}}

			ch, _ := New("https://example.com", Auth(dummyAuth.Builder()))
			ch.GetRSA("/example-rsa")
			url := dummyAuth.Request.URL
			Expect(url.String()).To(ContainSubstring("https://example.com/api/v1/data"))
			Expect(url.String()).To(ContainSubstring("name=%2Fexample-rsa"))
			Expect(url.String()).To(ContainSubstring("versions=1"))

			Expect(dummyAuth.Request.Method).To(Equal(http.MethodGet))
		})

		Context("when successful", func() {
			It("returns a rsa credential", func() {
				responseString := `{
				  "data": [
					{
					  "id": "67fc3def-bbfb-4953-83f8-4ab0682ad677",
					  "name": "/example-rsa",
					  "type": "rsa",
					  "value": {
						"public_key": "public-key",
						"private_key": "private-key"
					  },
					  "version_created_at": "2017-01-01T04:07:18Z"
					}
				  ]
				}`

				dummyAuth := &DummyAuth{Response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(responseString)),
				}}

				ch, _ := New("https://example.com", Auth(dummyAuth.Builder()))
				cred, err := ch.GetRSA("/example-rsa")
				Expect(err).ToNot(HaveOccurred())
				Expect(cred.Value).To(Equal(values.RSA{
					PublicKey:  "public-key",
					PrivateKey: "private-key",
				}))
			})
		})

		Context("when response body cannot be unmarshalled", func() {
			It("returns an error", func() {
				dummyAuth := &DummyAuth{Response: &http.Response{
					Body: ioutil.NopCloser(bytes.NewBufferString("something-invalid")),
				}}
				ch, _ := New("https://example.com", Auth(dummyAuth.Builder()))
				_, err := ch.GetRSA("/example-cred")

				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("GetSSH()", func() {
		It("requests the credential by name", func() {
			dummyAuth := &DummyAuth{Response: &http.Response{
				Body: ioutil.NopCloser(bytes.NewBufferString("")),
			}}

			ch, _ := New("https://example.com", Auth(dummyAuth.Builder()))
			ch.GetSSH("/example-ssh")
			url := dummyAuth.Request.URL
			Expect(url.String()).To(ContainSubstring("https://example.com/api/v1/data"))
			Expect(url.String()).To(ContainSubstring("name=%2Fexample-ssh"))
			Expect(url.String()).To(ContainSubstring("versions=1"))

			Expect(dummyAuth.Request.Method).To(Equal(http.MethodGet))
		})

		Context("when successful", func() {
			It("returns a ssh credential", func() {
				responseString := `{
				  "data": [
					{
					  "id": "some-id",
					  "name": "/example-ssh",
					  "type": "ssh",
					  "value": {
						"public_key": "public-key",
						"private_key": "private-key",
						"public_key_fingerprint": "public-key-fingerprint"
					  },
					  "version_created_at": "2017-01-01T04:07:18Z"
					}
				  ]
				}`

				dummyAuth := &DummyAuth{Response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(responseString)),
				}}

				ch, _ := New("https://example.com", Auth(dummyAuth.Builder()))
				cred, err := ch.GetSSH("/example-ssh")
				Expect(err).ToNot(HaveOccurred())
				Expect(cred.Value.PublicKeyFingerprint).To(Equal("public-key-fingerprint"))
				Expect(cred.Value.SSH).To(Equal(values.SSH{
					PublicKey:  "public-key",
					PrivateKey: "private-key",
				}))
			})
		})

		Context("when response body cannot be unmarshalled", func() {
			It("returns an error", func() {
				dummyAuth := &DummyAuth{Response: &http.Response{
					Body: ioutil.NopCloser(bytes.NewBufferString("something-invalid")),
				}}
				ch, _ := New("https://example.com", Auth(dummyAuth.Builder()))
				_, err := ch.GetSSH("/example-cred")

				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("GetJSON()", func() {
		It("requests the credential by name", func() {
			dummyAuth := &DummyAuth{Response: &http.Response{
				Body: ioutil.NopCloser(bytes.NewBufferString("")),
			}}

			ch, _ := New("https://example.com", Auth(dummyAuth.Builder()))
			ch.GetJSON("/example-json")
			url := dummyAuth.Request.URL
			Expect(url.String()).To(ContainSubstring("https://example.com/api/v1/data"))
			Expect(url.String()).To(ContainSubstring("name=%2Fexample-json"))
			Expect(url.String()).To(ContainSubstring("versions=1"))

			Expect(dummyAuth.Request.Method).To(Equal(http.MethodGet))
		})

		Context("when successful", func() {
			It("returns a json credential", func() {
				responseString := `{
				  "data": [
					{
					  "id": "some-id",
					  "name": "/example-json",
					  "type": "json",
					  "value": {
						"key": 123,
						"key_list": [
						  "val1",
						  "val2"
						],
						"is_true": true
					  },
					  "version_created_at": "2017-01-01T04:07:18Z"
					}
				  ]
				}`

				dummyAuth := &DummyAuth{Response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(responseString)),
				}}

				ch, _ := New("https://example.com", Auth(dummyAuth.Builder()))
				cred, err := ch.GetJSON("/example-json")
				Expect(err).ToNot(HaveOccurred())
				Expect([]byte(cred.Value)).To(MatchJSON(`{
						"key": 123,
						"key_list": [
						  "val1",
						  "val2"
						],
						"is_true": true
					}`))
			})
		})

		Context("when response body cannot be unmarshalled", func() {
			It("returns an error", func() {
				dummyAuth := &DummyAuth{Response: &http.Response{
					Body: ioutil.NopCloser(bytes.NewBufferString("something-invalid")),
				}}
				ch, _ := New("https://example.com", Auth(dummyAuth.Builder()))
				_, err := ch.GetJSON("/example-cred")

				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("GetValue()", func() {
		It("requests the credential by name", func() {
			dummyAuth := &DummyAuth{Response: &http.Response{
				Body: ioutil.NopCloser(bytes.NewBufferString("")),
			}}

			ch, _ := New("https://example.com", Auth(dummyAuth.Builder()))
			ch.GetValue("/example-value")
			url := dummyAuth.Request.URL
			Expect(url.String()).To(ContainSubstring("https://example.com/api/v1/data"))
			Expect(url.String()).To(ContainSubstring("name=%2Fexample-value"))
			Expect(url.String()).To(ContainSubstring("versions=1"))

			Expect(dummyAuth.Request.Method).To(Equal(http.MethodGet))
		})

		Context("when successful", func() {
			It("returns a value credential", func() {
				responseString := `{
				  "data": [
					{
					  "id": "some-id",
					  "name": "/example-value",
					  "type": "value",
					  "value": "some-value",
					  "version_created_at": "2017-01-05T01:01:01Z"
				}]}`

				dummyAuth := &DummyAuth{Response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(responseString)),
				}}

				ch, _ := New("https://example.com", Auth(dummyAuth.Builder()))
				cred, err := ch.GetValue("/example-value")
				Expect(err).ToNot(HaveOccurred())
				Expect(cred.Value).To(Equal(values.Value("some-value")))
			})
		})

		Context("when response body cannot be unmarshalled", func() {
			It("returns an error", func() {
				dummyAuth := &DummyAuth{Response: &http.Response{
					Body: ioutil.NopCloser(bytes.NewBufferString("something-invalid")),
				}}
				ch, _ := New("https://example.com", Auth(dummyAuth.Builder()))
				_, err := ch.GetValue("/example-cred")

				Expect(err).To(HaveOccurred())
			})
		})
	})

	DescribeTable("request fails due to network error",
		func(performAction func(*CredHub) error) {
			networkError := errors.New("Network error occurred")
			dummyAuth := &DummyAuth{Error: networkError}
			ch, _ := New("https://example.com", Auth(dummyAuth.Builder()))

			err := performAction(ch)

			Expect(err).To(Equal(networkError))
		},

		Entry("GetNVersions", func(ch *CredHub) error {
			_, err := ch.GetNVersions("/example-password", 47)
			return err
		}),
		Entry("GetLatestVersion", func(ch *CredHub) error {
			_, err := ch.GetLatestVersion("/example-password")
			return err
		}),
		Entry("GetPassword", func(ch *CredHub) error {
			_, err := ch.GetPassword("/example-password")
			return err
		}),
		Entry("GetCertificate", func(ch *CredHub) error {
			_, err := ch.GetCertificate("/example-certificate")
			return err
		}),
		Entry("GetUser", func(ch *CredHub) error {
			_, err := ch.GetUser("/example-password")
			return err
		}),
		Entry("GetRSA", func(ch *CredHub) error {
			_, err := ch.GetRSA("/example-password")
			return err
		}),
		Entry("GetSSH", func(ch *CredHub) error {
			_, err := ch.GetSSH("/example-password")
			return err
		}),
		Entry("GetJSON", func(ch *CredHub) error {
			_, err := ch.GetJSON("/example-password")
			return err
		}),
		Entry("GetValue", func(ch *CredHub) error {
			_, err := ch.GetValue("/example-password")
			return err
		}),
	)
})
