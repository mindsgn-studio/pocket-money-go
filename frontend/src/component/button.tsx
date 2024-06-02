import { Button as DefaultButton, Text } from '@chakra-ui/react'

function Button(props: any) {
    console.log(props)
    return (
        <DefaultButton
            onClick={props.onClick}
            backgroundColor='black'>
            <Text
                color="white"
                fontFamily={"SF Rounded Bold"}>
                    {props.text}
            </Text>
        </DefaultButton>
    )
}

export default Button
