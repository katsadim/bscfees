import * as React from "react";
import {FunctionComponent} from "react";
import {Card, CardActions, CardContent, Divider, makeStyles, Typography} from "@material-ui/core";
import {DetailsPopover} from "./DetailsPopover";

type Props = {
    bnbusdPrice: number
    ethusdPrice: number
    bscFees: number
    ethFees: number
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
                <Typography variant="h6" component="h2">
                    Bsc Fees: {props.bscFees.toFixed(props.numOfDecimals)}BNB
                </Typography>
                <Typography variant="h6" component="h2">
                    Bsc Fees: {(props.bnbusdPrice * props.bscFees).toFixed(props.numOfDecimals)}$
                </Typography>
                <Divider/>
                <Typography variant="h6" component="h2">
                    Eth Fees: {props.ethFees.toFixed(props.numOfDecimals)}ETH
                </Typography>
                <Typography variant="h6" component="h2">
                    Eth Fees: {(props.ethusdPrice * props.ethFees).toFixed(props.numOfDecimals)}$
                </Typography>
                <Divider/>
                <Typography variant="h5" component="h2">
                    Grand Total: {(
                    (props.bnbusdPrice * props.bscFees) +
                    (props.ethusdPrice * props.ethFees)
                ).toFixed(props.numOfDecimals)}$
                </Typography>
            </CardContent>
            <CardActions>
                <DetailsPopover bnbusdPrice={props.bnbusdPrice} ethusdPrice={props.ethusdPrice}/>
            </CardActions>

        </Card>
    );
}
