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
} from "@chakra-ui/react";
import { EditIcon, DeleteIcon } from "@chakra-ui/icons";
import { Check } from "../model";
import { useHistory, Link as RouterLink } from "react-router-dom";

interface ChecksListProps {
  checks: Check[];
  onDelete(id: string): void;
}

const ChecksTable = ({ checks, onDelete }: ChecksListProps) => {
  const history = useHistory();
  console.log("CHECKS", checks);

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
                  {/* onClick={() => {
                    history.push("/check", check);
                  }} */}
                  {/* > */}
                  {check.name}
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
                <Switch isChecked={check.active} />
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
