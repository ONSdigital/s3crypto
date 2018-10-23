package main

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"io"
	"io/ioutil"
	"os"

	"github.com/ONSdigital/go-ns/log"
	"github.com/ONSdigital/s3crypto"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func main() {

	f, err := ioutil.ReadFile("testdata/private.pem")
	if err != nil {
		panic(err)
	}

	block, _ := pem.Decode(f)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		panic(err)
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}

	region := "eu-west-1"

	sess, err := session.NewSession(&aws.Config{Region: aws.String(region)})
	if err != nil {
		panic(err)
	}

	svc := s3crypto.NewUploader(sess, &s3crypto.Config{PrivateKey: privateKey})

	bucket := "dp-frontend-florence-file-uploads"
	key := "cpicoicoptest.csv"

	b, err := ioutil.ReadFile("testdata/" + key)
	if err != nil {
		log.Error(err, nil)
		return
	}

	acl := "public-read"
	input := &s3manager.UploadInput{
		Body:   bytes.NewReader(b),
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}
	input.ACL = &acl

	cr, err := svc.Upload(input)
	if err != nil {
		log.Error(err, nil)
		return
	}

	log.Info("upload completed", log.Data{"result": cr})

	log.Info("now getting file...", nil)

	svcDown := s3crypto.New(sess, &s3crypto.Config{PrivateKey: privateKey})

	getInput := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	out, err := svcDown.GetObject(getInput)
	if err != nil {
		log.Error(err, nil)
		return
	}

	newf, err := os.Create("newcpi.csv")
	if err != nil {
		log.Error(err, nil)
		return
	}
	defer newf.Close()

	_, err = io.Copy(newf, out.Body)
	if err != nil {
		log.Error(err, nil)
		return
	}

}
