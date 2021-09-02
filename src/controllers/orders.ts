import { Request, Response } from 'express';
import { sendNewOrderMsg } from '../ds';
import { createOrder } from '../services/orders';

export const addOrder = async (req: Request, res: Response) => {
  //TODO add validation
  const order = req.body as Order;
  try {
    const createdOrder = await createOrder({ ...order, userID: req.userID });

    await sendNewOrderMsg(createdOrder);

    res.json(createdOrder);
  } catch (err) {
    res.send(err).status(500);
  }
};
