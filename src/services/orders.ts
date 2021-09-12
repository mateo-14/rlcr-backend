import { customAlphabet } from 'nanoid';
import { OrderMode, OrderStatus } from '../../@types';
import firestore from '../firestore';
import settingsService from '../services/settings';

const nanoid = customAlphabet(process.env.NANOID_ALPHABET!, 10);

const mapOrderSnapshotToOrder = (
  orderSnap: FirebaseFirestore.DocumentSnapshot<FirebaseFirestore.DocumentData>
): Order => ({
  ...(orderSnap.data() as Order),
  id: orderSnap.id,
  createdAt: orderSnap.data()?.createdAt.toMillis(),
});

export const createOrder = async (order: Order): Promise<Order> => {
  const settings = await settingsService.getSettings();
  order = {
    ...order,
    createdAt: Date.now(),
    price: (order.mode === OrderMode.Buy ? settings.creditBuyValue : settings.creditSellValue) * order.credits,
    status: OrderStatus.Pending,
  };
  const id = nanoid();
  return firestore
    .collection('users')
    .doc(order.userID)
    .collection('orders')
    .doc(id)
    .create({ ...order, createdAt: new Date(order.createdAt) })
    .then(() => {
      return { ...order, id };
    });
};

export const getOrders = async (userID: string, startAt?: string): Promise<Order[]> => {
  const orders = firestore.collection('users').doc(userID).collection('orders');
  const orderRef = startAt && (await orders.doc(startAt).get());
  let query = orders.orderBy('createdAt', 'desc').limit(10);

  if (orderRef && orderRef.exists) {
    query = query.startAfter(orderRef);
  }

  return query.get().then((orders) => orders.docs.map((docSnap) => mapOrderSnapshotToOrder(docSnap)));
};

export const getOrder = (userID: string, orderID: string): Promise<Order> => {
  return firestore
    .collection('users')
    .doc(userID)
    .collection('orders')
    .doc(orderID)
    .get()
    .then((docSnap) => {
      if (docSnap.exists && docSnap.data()) {
        return mapOrderSnapshotToOrder(docSnap);
      } else {
        throw new Error('Order does not exists');
      }
    });
};

export interface GetAllOrdersQuery {
  orderBy?: string;
  order?: FirebaseFirestore.OrderByDirection;
  status?: OrderStatus[];
  users?: string[];
  userID?: string;
}
export const getAll = (query: GetAllOrdersQuery): Promise<Order[]> => {
  let fsQuery: FirebaseFirestore.Query = query.userID
    ? firestore.collection('users').doc(query.userID).collection('orders')
    : firestore.collectionGroup('orders');

  if (query.orderBy) fsQuery = fsQuery.orderBy(query.orderBy, query.order);
  else fsQuery = fsQuery.orderBy('createdAt', query.order || 'desc');
  if (query.status) fsQuery = fsQuery.where('status', 'in', query.status);
  if (query.users) fsQuery = fsQuery.where('userID', 'in', query.users);

  return fsQuery.get().then((querySnap) => querySnap.docs.map((docSnap) => mapOrderSnapshotToOrder(docSnap)));
};
