import { Button as DefaultButton, Text } from '@chakra-ui/react'

function Button(props: any) {
    return (
        <DefaultButton
            backgroundColor='black'>
            <Text color="white">{props.text}</Text>
        </DefaultButton>
    )
}

export default Button
