import {
  Button,
  Center,
  Flex,
  Spinner,
  Stack,
  toast,
  useToast,
} from "@chakra-ui/react";
import React, { useEffect } from "react";
import { QueryCache, useQuery } from "react-query";
import { useHistory, useLocation, useParams } from "react-router-dom";
import { HISTORY_QUERY_KEY, QUERY_KEY } from "../constants";
import * as api from "../api/checks";
import Block from "../components/Container";
import StatusDetails from "../components/Status";
import { ArrowBackIcon } from "@chakra-ui/icons";

interface ParamTypes {
  id: string;
}

const CheckDetails = () => {
  const toast = useToast();
  const routerHistory = useHistory();
  const { id } = useParams<ParamTypes>();
  const { isSuccess, isLoading, isError, data: check } = useQuery(
    [QUERY_KEY, id],
    async () => api.getCheck(id)
  );

  const { isError: err2, data: history } = useQuery(
    [HISTORY_QUERY_KEY, id],
    async () => api.getHistory(id)
  );

  if (isLoading) {
    return <Spinner size="xl" />;
  }

  if (isError) {
    toast({
      position: "bottom-right",
      title: "Error",
      description: `There was an error fetching details for check ${id}`,
      status: "error",
      duration: 10000,
      isClosable: true,
    });
    routerHistory.goBack();
  }

  if (!isSuccess) {
    return <div>Error</div>;
  }

  const toolbar = (
    <Button
      onClick={() => routerHistory.goBack()}
      leftIcon={<ArrowBackIcon />}
      colorScheme="teal"
      variant="outline"
    >
      Go back
    </Button>
  );

  return (
    <>
      <Block title={check!.name} toolbar={toolbar}>
        <Stack spacing={4} w="100%">
          {/* {isLoading && <Spinner size="xl" />} */}
          {isSuccess &&
            history?.map((status) => <StatusDetails status={status} />)}
        </Stack>
      </Block>
    </>
  );
};

export default CheckDetails;