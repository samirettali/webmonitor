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
import { Check } from "../check";

interface ChecksListProps {
  checks: Check[];
  onDelete(id: string): void;
}

const ChecksTable = ({ checks, onDelete }: ChecksListProps) => {
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
          {checks.map(({ name, id, url, interval, active, email }: Check) => (
            <Tr key={id}>
              <Td>{name}</Td>
              <Td>
                <Link href={url} color="blue.500" isExternal>
                  {url}
                </Link>
              </Td>
              <Td isNumeric>{interval}</Td>
              <Td>{email}</Td>
              <Td>
                <Switch isChecked={active} />
              </Td>
              <Td>
                <Flex justify="space-evenly">
                  <EditIcon
                    boxSize={6}
                    _hover={{
                      color: "teal.500",
                      cursor: "pointer",
                    }}
                    onClick={() => onDelete(id!)}
                  />
                  <DeleteIcon
                    boxSize={6}
                    _hover={{
                      color: "teal.500",
                      cursor: "pointer",
                    }}
                    onClick={() => onDelete(id!)}
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
