import {documentGetInitialProps, DocumentHeadTags, DocumentHeadTagsProps,} from '@mui/material-nextjs/v15-pagesRouter';
import {DocumentProps, Head, Html, Main, NextScript } from 'next/document';
// or `v1X-pagesRouter` if you are using Next.js v1X

export default function MyDocument( props: DocumentProps & DocumentHeadTagsProps) {
    return (
        <Html lang="en">
            <Head>
                <DocumentHeadTags {...props} />
                ...
            </Head>
            <body>
            <Main/>
            <NextScript/>
            </body>
        </Html>
    );
}

MyDocument.getInitialProps = async (ctx: any) => {
    const finalProps = await documentGetInitialProps(ctx);
    return finalProps;
};
