import React from "react";
import { Box, Button, IconButton, useColorMode } from "@chakra-ui/react";
import { MoonIcon, SunIcon } from "@chakra-ui/icons";

const Toggler = () => {
  const { colorMode, toggleColorMode } = useColorMode();
  return (
    <IconButton
      aria-label="Toggle color mode"
      onClick={toggleColorMode}
      icon={colorMode === "light" ? <SunIcon /> : <MoonIcon />}
    />
    // <Box colorScheme="teal" onClick={toggleColorMode}>
    //   {colorMode === "light" ? <SunIcon /> : <MoonIcon />}
    // </Box>
  );
};

export default Toggler;
