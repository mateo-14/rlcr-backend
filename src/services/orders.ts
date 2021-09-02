import { customAlphabet } from 'nanoid';
import firestore from '../firestore';

const nanoid = customAlphabet(process.env.NANOID_ALPHABET!, 10);

export const createOrder = (order: Order): Promise<Order> => {
  const id = nanoid();
  return firestore
    .collection('users')
    .doc(order.userID)
    .collection('orders')
    .doc(id)
    .create(order)
    .then(() => ({ ...order, id, createdAt: new Date() }));
};
