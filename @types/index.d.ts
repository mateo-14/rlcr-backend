declare namespace Express {
  interface Request {
    userID: string;
  }
}

interface Order {
  id: string;
  userID: string;
  credits: string;
  paymentMethodId: string;
  mode: number;
  dni: string;
  account: string;
  paymentAccount?: string;
  cvuAlias: string;
  createdAt: Date;
}
