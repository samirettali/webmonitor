import React from "react";
import {
  Table,
  Thead,
  Link,
  Center,
  Tbody,
  Tr,
  Th,
  Td,
  Switch,
  Flex,
  toast,
  useToast,
} from "@chakra-ui/react";
import { EditIcon, DeleteIcon } from "@chakra-ui/icons";
import { Check, CheckUpdate } from "../model";
import { useHistory, Link as RouterLink } from "react-router-dom";
import * as api from "../api/checks";
import { QUERY_KEY } from "../constants";
import { QueryClient, useMutation, useQueryClient } from "react-query";

interface ChecksListProps {
  checks: Check[];
  onDelete(id: string): void;
}

interface ToggleMutationData {
  id: string;
  upd: CheckUpdate;
}

const ChecksTable = ({ checks, onDelete }: ChecksListProps) => {
  const history = useHistory();
  const toast = useToast();
  const queryClient = useQueryClient();

  const toggleMutation = useMutation(
    (data: ToggleMutationData) => api.updateCheck(data.id, data.upd),
    {
      onMutate: async ({ id, upd }: ToggleMutationData) => {
        await queryClient.cancelQueries(QUERY_KEY);
        const previousChecks = queryClient.getQueryData<Check[]>(QUERY_KEY);
        if (previousChecks !== undefined) {
          queryClient.setQueryData<Check[]>(
            QUERY_KEY,
            previousChecks.map((check: Check) => {
              if (check.id != id) return check;
              else return { ...check, active: upd.active! };
            })
          );
        }
        return { previousChecks };
      },
      onSuccess: () => {
        queryClient.invalidateQueries(QUERY_KEY);
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
          queryClient.setQueryData<Check[]>(QUERY_KEY, context.previousChecks);
        }
      },
    }
  );

  const toggleCheck = ({ id, active }: Check) => {
    const upd: CheckUpdate = { active: !active };
    const data: ToggleMutationData = { id: id!, upd };
    toggleMutation.mutate(data);
  };

  return (
    <>
      <Table variant="simple">
        <Thead>
          <Tr>
            <Th>Name</Th>
            <Th>URL</Th>
            <Th isNumeric>Interval</Th>
            <Th>Email</Th>
            <Th>Active</Th>
            <Th>
              <Center>Actions</Center>
            </Th>
          </Tr>
        </Thead>
        <Tbody>
          {checks.map((check: Check) => (
            <Tr key={check.id}>
              <Td>
                <RouterLink
                  to={{ pathname: `/check/${check.id}`, state: check }}
                >
                  <Link color="blue.500">{check.name}</Link>
                </RouterLink>
              </Td>
              <Td>
                <Link href={check.url} color="blue.500" isExternal>
                  {check.url}
                </Link>
              </Td>
              <Td isNumeric>{check.interval}</Td>
              <Td>{check.email}</Td>
              <Td>
                <Switch
                  isChecked={check.active}
                  onChange={() => toggleCheck(check)}
                />
              </Td>
              <Td>
                <Flex justify="space-evenly">
                  <EditIcon
                    boxSize={6}
                    _hover={{
                      color: "teal.500",
                      cursor: "pointer",
                    }}
                    onClick={() => onDelete(check.id!)}
                  />
                  <DeleteIcon
                    boxSize={6}
                    _hover={{
                      color: "teal.500",
                      cursor: "pointer",
                    }}
                    onClick={() => onDelete(check.id!)}
                  />
                </Flex>
              </Td>
            </Tr>
          ))}
        </Tbody>
      </Table>
    </>
  );
};

export default ChecksTable;
