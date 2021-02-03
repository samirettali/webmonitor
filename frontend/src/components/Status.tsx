import React from "react";
import { Box, Stack, Text } from "@chakra-ui/react";
import { format } from "date-fns";
import { Status } from "../model";

interface StatusProps {
  status: Status;
}

const StatusDetails: React.FC<StatusProps> = ({ status }) => {
  const { date, content } = status;
  const formattedDate = format(date, "dd-MM-yyyy HH:mm:ss");
  return (
    <Stack spacing={2} p={4} w="100%">
      <Box>
        <Text fontSize="xl">{formattedDate}</Text>
      </Box>
      <Box
        as="pre"
        overflow="auto"
        fontSize="sm"
        bg="gray.100"
        rounded="sm"
        p={4}
      >
        {content}
      </Box>
    </Stack>
  );
};

export default StatusDetails;
