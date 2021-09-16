import { Request, Response } from 'express';
import { validationResult } from 'express-validator';
import { sendNewOrderMsg } from '../ds';
import * as ordersService from '../services/orders';
import { GetAllOrdersQuery, getAll, update } from '../services/orders';

export async function addOrder(req: Request, res: Response) {
  const errors = validationResult(req);
  if (!errors.isEmpty()) {
    return res.status(400).json({ errors: errors.mapped() });
  }

  try {
    const order = req.body as Order;
    const createdOrder = await ordersService.createOrder({ ...order, userID: req.userID });
    await sendNewOrderMsg(createdOrder);

    res.json(createdOrder);
  } catch (err) {
    console.error(err);
    res.status(500).send(err);
  }
}

export function getOrders(req: Request, res: Response) {
  if (typeof req.query.startAfter === 'string' || !req.query.startAfter) {
    ordersService
      .getOrders(req.userID, req.query.startAfter)
      .then((orders) => res.json(orders))
      .catch((err) => res.send(err).status(500));
  } else {
    res.sendStatus(400);
  }
}

export function getOrder(req: Request, res: Response) {
  ordersService
    .getOrder(req.userID, req.params.id)
    .then((order) => {
      res.json(order);
    })
    .catch(() => {
      res.sendStatus(404);
    });
}

export function getAllOrders(req: Request, res: Response) {
  const query = req.query as GetAllOrdersQuery;
  if (query.status) query.status = query.status.map((status) => parseInt(status.toString()));

  getAll(query)
    .then((orders) => {
      res.json(orders);
    })
    .catch((err) => {
      console.error(err);
      res.sendStatus(500);
    });
}

export function updateOrder(req: Request, res: Response) {
  if (req.body.status && req.body.userID) {
    return update(req.body.userID, req.params.id, { status: req.body.status })
      .then((order) => res.json(order))
      .catch((err) => {
        console.error(err);
        res.sendStatus(404);
      });
  }
  res.sendStatus(400);
}
