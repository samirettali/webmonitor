import { Box, Center, Flex, Heading, Stack } from "@chakra-ui/react";
import React from "react";

interface ContainerProps {
  title: string;
  toolbar?: React.ReactNode;
  aside?: string;
}

const Block: React.FC<ContainerProps> = ({ title, toolbar, children }) => {
  return (
    <Box w="100%" minH={40} px={16}>
      <Stack bg="white" spacing={16} p={8} boxShadow="xl" rounded="lg">
        <Flex justify="space-between">
          <Box>
            <Heading as="h2">{title}</Heading>
          </Box>
          {toolbar && <Box>{toolbar}</Box>}
        </Flex>
        <Center>{children}</Center>
      </Stack>
    </Box>
  );
};

export default Block;
