import {
    Table,
    Thead,
    Tbody,
    Tfoot,
    Tr,
    Th,
    Td,
} from "@chakra-ui/react"
import { ApiInterface } from "./api_interface";

export const TableData = (props: ApiInterface) => {
    const { html_version, title, h1_count, h2_count, h3_count, h4_count,
        h5_count, h6_count, internal_links_count, external_links_count,
        inaccessible_links_count, has_login_form } = props;

    const table_body = [
        {key: "HTML version", val: html_version },
        {key: "Page Title", val: title },
        {key: "H1 count", val: h1_count},
        {key: "H2 count", val: h2_count},
        {key: "H3 count", val: h3_count},
        {key: "H4 count", val: h4_count},
        {key: "H5 count", val: h5_count},
        {key: "H6 count", val: h6_count},
        {key: "Internal Links Count", val: internal_links_count},
        {key: "External Links Count", val: external_links_count},
        {key: "Inaccessible Links Count", val: inaccessible_links_count},
        {key: "Has Login Form", val: has_login_form ? "Yes" : "No"}
    ]
    return (
        <Table variant="striped" colorScheme="gray" width="50vw" margin="auto" mt="36px" border="1px solid gray">
            <Thead>
                <Tr>
                    <Th>Page Property</Th>
                    <Th>Value</Th>
                </Tr>
            </Thead>
            <Tbody>
                {table_body.map(({key, val}, idx) => (
                    <Tr key={idx}>
                        <Td>{key}</Td>
                        <Td>{val}</Td>
                    </Tr>
                ))}
            </Tbody>
            <Tfoot>
                <Tr>
                    <Th>Page Property</Th>
                    <Th>Value</Th>
                   
                </Tr>
            </Tfoot>
        </Table>
    )
}