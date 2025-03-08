# DICOM API

This DICOM API is used to upload a DICOM file, extract and return any DICOM header attribute based on a DICOM Tag as a query parameter, and finally convert the uploaded file into a PNG for browser-based viewing using REST API.

## Uploading the DICOM File

To upload a DICOM file, run the following command:

```bash
curl -X POST -F "file=@filepath" http://localhost:8080/dicom/v1/upload
```

If the upload is successful, the command returns the following message:

```text
File uploaded successfully: file_name.dcm
```

Example:

```bash
curl -X POST -F "file=@testdata/client/multi_frame.dcm" http://localhost:8080/dicom/v1/upload
```

## Getting Metadata of a DICOM Tag

Use the `file_name.dcm` returned from the upload to get the metadata value of the desired tag. The command for getting the metadata of a given tag is as follows:

```bash
curl -X GET "http://localhost:8080/dicom/v1/metadata?file=file_name.dcm&tag=tag_numbers"
```

For more information about DICOM tags, see this page: [DICOM Tags](https://www.dicomlibrary.com/dicom/dicom-tags/).

If there is no error, the tag value will be outputted.

Example:

```bash
curl -X GET "http://localhost:8080/dicom/v1/metadata?file=multi_frame.dcm&tag=(0010,0010)"

```


## Converting DICOM to PNG

Using the `file_name.dcm` outputted after uploading the file, you can convert the images in the DICOM file to PNG files using the following command. 

```bash
curl -X GET "http://localhost:8080/dicom/v1/png/conversion?file=file_name.dcm." -o output.zip
```

The PNG files will be downloaded as a zip file, in this example, `output.zip`.

Example:

```bash
curl -X GET "http://localhost:8080/dicom/v1/png/conversion?file=multi_frame.dcm" -o output.zip
```
