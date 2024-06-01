import Button from "./button"
import { Box , Text, PinInput, PinInputField, HStack, Heading, Center} from "@chakra-ui/react"

function NewWalletForm() {
    return (
        <Box >
            <Box
                marginY={10}>
                <Heading
                    size={"sm"}>
                    Create New Password
                </Heading>

                <Center>
                <HStack
                    marginY={10}>
                    <PinInput>
                        <PinInputField />
                        <PinInputField />
                        <PinInputField />
                        <PinInputField />
                    </PinInput>
                </HStack>
                </Center>
               
            </Box>
            
            <Button text='Next' />
        </Box>
    )
}

export default NewWalletForm
