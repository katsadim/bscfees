import * as React from "react";
import {FunctionComponent} from "react";
import {Card, CardActions, CardContent, makeStyles, Typography} from "@material-ui/core";
import {DetailsPopover} from "./DetailsPopover";

type Props = {
    bnbusdPrice: number
    fees: number
    numOfDecimals: number
}

const useStyles = makeStyles({
    root: {
        minWidth: "100%",
    },
});

export const ResultCard: FunctionComponent<Props> = (props: Props) => {
    const classes = useStyles();

    return (
        <Card className={classes.root}>
            <CardContent>
                <Typography variant="h5" component="h2">
                    Fees in BNB: {props.fees.toFixed(props.numOfDecimals)}
                </Typography>
                <Typography variant="h5" component="h2">
                    Fees in USD: {(props.bnbusdPrice * props.fees).toFixed(props.numOfDecimals)}
                </Typography>
            </CardContent>
            <CardActions>
                <DetailsPopover bnbusdPrice={props.bnbusdPrice}/>
            </CardActions>

        </Card>
    );
}
