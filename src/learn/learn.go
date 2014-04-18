package learn

import (
    "lbfgs"
    "crfpp"
    "conf"
)

func Learn(template, train, model string) error {
    c := conf.GetConfig()
    freq := c.GetFreq()
    maxLimit := c.GetMaxLimit()
    cost := c.GetCost()
    eta := c.GetEta()

    gr := new(crfpp.Group)
    e := gr.ReadFile(template, train)
    if e!= nil {
        return e
    }
    e = gr.ProcessData(model, freq, eta, cost, maxlimit)
    if e!= nil {
        return e
    }
    return nil
}
