import { useState } from "react"
import Button from "./button"
import { Box, PinInput, PinInputField, HStack, Heading, Center} from "@chakra-ui/react"

function NewWalletForm() {
    const [step, setStep] = useState<number>(1)
    const [passcode, setPasscode] = useState<string>("")
    const [verify, setVerify] = useState<string>("")

    const next = () => {
        if(step == 1){

        }else if(step == 2){
            verifyCode()
        }

        setStep(step+1)
    }

    const verifyCode = () => {
        window.location.href = "/home"
    }

    return (
        <Box >
            {
                step == 1?
                <Box
                    marginY={10}>
                    <Heading
                        fontFamily={"SF Rounded Regular"}
                        size={"sm"}
                        color="#908C87">
                        Create New Passcode
                    </Heading>

                    <Center>
                        <HStack
                            marginY={2}>
                            <PinInput>
                                <PinInputField />
                                <PinInputField />
                                <PinInputField />
                                <PinInputField />
                            </PinInput>
                        </HStack>
                    </Center>
                </Box>
                :
                <Box
                    marginY={10}>
                    <Heading
                        fontFamily={"SF Rounded Regular"}
                        size={"sm"}
                        color="#908C87">
                       Verify Passcode
                    </Heading>

                    <Center>
                        <HStack
                            marginY={2}>
                            <PinInput
                                onChange={(input) => setVerify(input)}>
                                <PinInputField />
                                <PinInputField />
                                <PinInputField />
                                <PinInputField />
                            </PinInput>
                        </HStack>
                    </Center>
                </Box>
            }

            <Button text={step==1? 'Next': "Verify"} onClick={next}/>
        </Box>
    )
}

export default NewWalletForm
