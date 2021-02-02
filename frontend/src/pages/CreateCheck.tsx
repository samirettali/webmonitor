import React, { useEffect } from "react";
import { Formik, Form, Field } from "formik";
import {
  Input,
  Button,
  FormControl,
  FormLabel,
  FormErrorMessage,
  FormHelperText,
  Box,
  Center,
  Stack,
  Heading,
  Select,
  useToast,
  Spinner,
} from "@chakra-ui/react";
import * as Yup from "yup";
import { useHistory } from "react-router-dom";
import humanizeDuration from "humanize-duration";


import { Check } from "../check";
import { QueryCache, useMutation, useQueryClient } from "react-query";
import * as api from "../api/checks";
import axios from "axios";
import { QUERY_KEY, INTERVALS } from "../constants";

const CreateSchema = Yup.object().shape({
  name: Yup.string().min(3).max(30).required("Required"),
  url: Yup.string().required("Required"),
  // url: Yup.string().url("Invalid URL").required("Required"),
  interval: Yup.number()
    .min(1)
    .integer("Invalid interval")
    .required("Required"),
  email: Yup.string().email("Invalid Email").required("Required"),
});

function isError(error: unknown): error is Error {
  return error instanceof Error;
}

const CreateCheck = () => {
  const toast = useToast();
  const history = useHistory();
  const queryClient = useQueryClient();
  const queryCache = new QueryCache({
    onError: (error) => {
      console.log("QueryCache error: ", error);
    },
  });

  const initialValues: Check = {
    name: "My check",
    url: "http://localhost:8080/test",
    interval: 0,
    email: "mail@example.com",
    active: true,
  };


  const { mutate } = useMutation(api.createCheck, {
      onMutate: async (newCheck: Check) => {
        await queryClient.cancelQueries(QUERY_KEY);
        const previousChecks = queryClient.getQueryData<Check[]>(QUERY_KEY);
        // Temporary id
        newCheck.id = Date.now().toString();
        if (previousChecks) {
          queryClient.setQueryData<Check[]>(QUERY_KEY, [
            ...previousChecks,
            newCheck
          ]);
        } else {
          queryClient.setQueryData<Check[]>(QUERY_KEY, [
            newCheck
          ]);
        }
        return { previousChecks };
      },
      onError: (err, variables, context) => {
        toast({
          position: "bottom-right",
          title: "An error occurred",
          description: "There was an error",
          status: "error",
          duration: 10000,
          isClosable: true,
        }); 
        if (context?.previousChecks) {
          queryClient.setQueryData<Check[]>(QUERY_KEY, context.previousChecks)
        }
      },
      onSuccess: (newCheck) => {
        queryClient.invalidateQueries(QUERY_KEY);
        history.push("/dashboard");
      },
    }
  );

  return (
    <Stack spacing="1rem" bg="white" shadow="xl" p={8} rounded="xl">
      <Heading as="h2">Create a new check</Heading>
      <Center>
        <Box width="full">
          <Formik
            initialValues={initialValues}
            validationSchema={CreateSchema}
            onSubmit={(values: Check, { setSubmitting }) => {
              const check: Check = {
                ...values,
                interval: typeof values.interval === "string" ? parseInt(values.interval, 10) : values.interval
              }
              // values.interval = parseInt(values.interval, 10);
              mutate(check);
              setSubmitting(false);
            }}
          >
            {({ errors, touched, isSubmitting }) => (
              <Form>
              <Stack spacing={4}>
                <Field name="name">
                  {({ field, form }: any) => (
                    <FormControl
                      isInvalid={form.errors.name && form.touched.name}
                    >
                      <FormLabel htmlFor="url">Name</FormLabel>
                      <Input {...field} id="name" placeholder="My check" />
                      <FormErrorMessage>{form.errors.name}</FormErrorMessage>
                    </FormControl>
                  )}
                </Field>
                <Field name="url">
                  {({ field, form }: any) => (
                    <FormControl
                      isInvalid={form.errors.url && form.touched.url}
                    >
                      <FormLabel htmlFor="url">URL</FormLabel>
                      <Input
                        {...field}
                        id="url"
                        placeholder="https://example.com"
                      />
                      <FormErrorMessage>{form.errors.url}</FormErrorMessage>
                    </FormControl>
                  )}
                </Field>
                <Field name="interval">
                  {({ field, form }: any) => (
                    <FormControl
                      isInvalid={form.errors.interval && form.touched.interval}
                    >
                      <FormLabel htmlFor="interval">Interval</FormLabel>
                      <Select placeholder="Select an interval" {...field}>
                        {INTERVALS.map(interval => (
                          <option value={interval}>{humanizeDuration(interval * 1000)}</option>)
                        )}
                      </Select>
                      <FormErrorMessage>
                        {form.errors.interval}
                      </FormErrorMessage>
                    </FormControl>
                  )}
                </Field>
                <Field name="email">
                  {({ field, form }: any) => (
                    <FormControl
                      isInvalid={form.errors.email && form.touched.email}
                    >
                      <FormLabel htmlFor="email">Email</FormLabel>
                      <Input
                        {...field}
                        id="email"
                        placeholder="me@example.com"
                      />
                      <FormErrorMessage>{form.errors.email}</FormErrorMessage>
                    </FormControl>
                  )}
                </Field>
                <Button
                  mt={4}
                  colorScheme="teal"
                  isLoading={isSubmitting}
                  type="submit"
                  width="full"
                >
                Create
                </Button>
              </Stack>
              </Form>
            )}
          </Formik>
        </Box>
      </Center>
    </Stack>
  );
};

export default CreateCheck;
