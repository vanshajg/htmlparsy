import { useState } from 'react'
import './App.css';
import axios from 'axios';
import { Input, Select, Flex, Checkbox, Button, Tooltip, FormControl, FormErrorMessage, Heading } from '@chakra-ui/react'
import { AiOutlineInfoCircle } from 'react-icons/ai'
import { Form, Formik, Field } from 'formik'
import { TableData } from './Table'

const validateUrl = (url: string): string => {
  const regex = /(www\.)?[-a-zA-Z0-9@:%._\+~#=]{2,256}\.[a-z]{2,6}\b([-a-zA-Z0-9@:%_\+.~#?&//=]*)/
  if (url === "")
    return "url cannot be empty"
  if (!regex.test(url))
    return "url doesn't seem right"
  return ""
}

const API_STATUS = {
  "INIT": "INIT",
  "FETCHING": "FETCHING",
  "FETCHED": "FETCHED",
  "FAILED": "FAILED",
}

const App = () => {
  const [api_status, setApiStatus] = useState(API_STATUS.INIT)
  const [api_data, setApiData] = useState({
    "html_version": "",
    "title": "",
    "h1_count": 0,
    "h2_count": 0,
    "h3_count": 0,
    "h4_count": 0,
    "h5_count": 0,
    "h6_count": 0,
    "internal_links_count": 0,
    "external_links_count": 0,
    "inaccessible_links_count": 0,
    "has_login_form": false
  })
  return (
    <div className="App">
      <Heading textAlign="center" mb="8vh">HTML Parser</Heading>
      <Formik initialValues={{ url: "", protocol: "https", override_cache: false }} onSubmit={(values, actions) => {
        setApiStatus(API_STATUS.FETCHING);
        axios.get("/api/data", {
          params: {
            url: `${values.protocol}://${values.url}`,
            ignoreCache: values.override_cache
          }
        }).then(data => {
          setApiStatus(API_STATUS.FETCHED);
          setApiData(data.data);
        }).catch(err => {
          setApiStatus(API_STATUS.FAILED);
          actions.setFieldError("url", "Request for this URL failed, check the URL or try again");
        }).finally(() => {
          actions.setSubmitting(false)
        })
      }}>
        {(props) => (
          <Form style={{ display: "flex", justifyContent: "center", width: "100vw", height: "max-content", alignItems: "start" }}>
            <Field name="protocol">
              {({ field, form }: any) => (
                <FormControl width="max-content">
                  <Select size="lg" variant="filled" defaultValue="https" width="max-content" id="protocol" onChange={props.handleChange}>
                    <option value="http">http://</option>
                    <option value="https">https://</option>
                  </Select>
                </FormControl>
              )}
            </Field>
            <Field name="url" validate={validateUrl}>
              {({ field, form }: any) => (
                <FormControl isRequired width="20vw" isInvalid={form.errors.url && form.touched.url}>
                  <Input size="lg" {...field} id="url" placeholder="google.com" />
                  <FormErrorMessage>{form.errors.url}</FormErrorMessage>
                </FormControl>
              )}
            </Field>
            <Button size="lg" colorScheme="blue" ml="4" type="submit" isLoading={props.isSubmitting}>Submit</Button>
            <Field>
              {() => (
                <FormControl width="max-content">
                  <Checkbox size="lg" ml="16" onChange={props.handleChange} id="override_cache">
                    <Tooltip label="URL data and computation are cached for ~6hours, check this tickbox to flush cache for this url and recompute" placement="right" openDelay={500}>
                      <Flex alignItems="center">
                        Override cache?
                       <AiOutlineInfoCircle style={{ marginLeft: "4px" }} />
                      </Flex>
                    </Tooltip>
                  </Checkbox>
                </FormControl>
              )}
            </Field>
          </Form>
        )}
      </Formik>
      {api_status === API_STATUS.FETCHED ? <TableData {...api_data} /> : null}
    </div>
  );
}

export default App;
