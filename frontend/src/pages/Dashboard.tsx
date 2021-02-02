import React, { useState, useEffect } from "react";
import {
  Container,
  Heading,
  Box,
  Spacer,
  Flex,
  Link,
  useToast,
  Spinner,
  Stack,
  Button,
  Center,
} from "@chakra-ui/react";
import { useQuery, useMutation, useQueryClient } from "react-query";
import { AddIcon } from "@chakra-ui/icons";
import { Link as RouterLink } from "react-router-dom";

import * as api from "../api/checks";
import ChecksTable from "../components/ChecksList";
import { QUERY_KEY } from "../constants";
import { Check } from "../check";

const Dashboard = () => {
  const toast = useToast();
  const queryClient = useQueryClient();
  const { isSuccess, isLoading, isError, data: checks } = useQuery(
    QUERY_KEY,
    api.getChecks
  );


  const deleteMutation = useMutation(api.deleteCheck, {
    onMutate: async (id: string) => {
      await queryClient.cancelQueries(QUERY_KEY);
      const previousChecks = queryClient.getQueryData(QUERY_KEY);
      queryClient.setQueryData(
        QUERY_KEY,
        checks.filter((check: Check) => check.id != id)
      );
      return { previousChecks };
    },
    onError: (err, id, context) => {
      const newValue = context ? context.previousChecks : [];
      queryClient.setQueryData(QUERY_KEY, newValue);
    },
    onSettled: () => {
      queryClient.invalidateQueries(QUERY_KEY);
    },
  });

  const onDelete = (id: string) => {
    deleteMutation.mutate(id);
  };

  return (
    <>
      <Box w="100%" minH={40} px={16}>
        <Stack bg="white" spacing={16} p={8} boxShadow="xl" rounded="lg">
          <Flex justify="space-between">
            <Box>
              <Heading as="h2">Your checks</Heading>
            </Box>
            <RouterLink to="/new">
              <Button colorScheme="teal">Create</Button>
            </RouterLink>
          </Flex>
          <Center>
            {isLoading && <Spinner size="xl" />}
            {isError && <Box>Could not connect to server.</Box>}
            {isSuccess && checks && (
              <ChecksTable checks={checks} onDelete={onDelete} />
            )}
          </Center>
        </Stack>
      </Box>
    </>
  );
};

export default Dashboard;
