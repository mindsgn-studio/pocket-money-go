import NewWalletForm from "../component/form"
import { useSteps, Box, StepSeparator, StepTitle, StepDescription, Stepper, Step, StepIndicator, StepStatus, StepIcon,StepNumber, Heading } from "@chakra-ui/react"

const steps = [
    { title: 'First', description: 'Contact Info' },
    { title: 'Second', description: 'Date & Time' },
    { title: 'Third', description: 'Select Rooms' },
]

function Onbaording() {
    const { activeStep } = useSteps({
        index: 1,
        count: steps.length,
      })

    return (
        <div id="app-home">
            <Box
                display={"flex"}
                flexDir={"column"}
                background="white"
                padding={20}
                borderRadius={10}>
                <Box marginY={10}>
                <Box
                    marginY={10}>
                    <Heading>Pocket Money</Heading>
                </Box>
                <Stepper index={activeStep}>
                    {steps.map((step, index) => (
                        <Step key={index}>
                        <StepIndicator>
                            <StepStatus
                            complete={<StepIcon />}
                            incomplete={<StepNumber />}
                            active={<StepNumber />}
                            />
                        </StepIndicator>
                        <StepSeparator />
                        </Step>
                    ))}
                </Stepper>
                </Box>
                
                <NewWalletForm />
            </Box>
        </div>
    )
}

export default Onbaording
