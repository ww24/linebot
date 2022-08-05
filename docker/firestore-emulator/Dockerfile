ARG BASE_IMAGE_TAG=406.0.0-emulators
FROM gcr.io/google.com/cloudsdktool/cloud-sdk:$BASE_IMAGE_TAG

ENV HOST_PORT=":8080"

ENTRYPOINT gcloud beta emulators firestore start --host-port=$HOST_PORT
