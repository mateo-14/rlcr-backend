import { Router } from 'express';
import { checkSchema, Schema } from 'express-validator';
import { OrderMode, PaymentMethodID } from '../../@types';
import { addOrder, getOrder, getOrders, getAllOrders } from '../controllers/orders';
import verifyToken from '../middlewares/verifyToken';
import { sanitizeCredits } from '../services/settings';
import validator from 'validator';
import isAdmin from '../middlewares/admin';
const router = Router();
const addOrderSchema: Schema = {
  paymentMethodID: {
    in: ['body'],
    customSanitizer: {
      options: (value) => {
        if (!value || !(value in PaymentMethodID)) {
          return PaymentMethodID.MercadoPago;
        }
        return value;
      },
    },
  },
  mode: {
    in: ['body'],
    customSanitizer: {
      options: (value) => {
        if (!value || !(value in OrderMode)) {
          return OrderMode.Buy;
        }
        return value;
      },
    },
  },
  credits: {
    in: ['body'],
    customSanitizer: {
      options: (value, { req }) => {
        value = parseInt(value);
        if (!isNaN(value)) {
          return sanitizeCredits(value, req.body.mode);
        }
      },
    },
    notEmpty: {
      errorMessage: 'Ingresá los créditos',
    },
  },
  account: {
    in: ['body'],
    isLength: { options: { min: 3, max: 80 } },
    errorMessage: 'Ingresá una cuenta válida',
  },
  dni: {
    in: ['body'],
    customSanitizer: {
      options: (value) => {
        if (value) return parseInt(value);
      },
    },
    custom: {
      options: (value, { req }) => {
        if (req.body.paymentMethodID !== PaymentMethodID.Transferencia) {
          return true;
        } else if (value && value > 1000000 && value < 100000000) {
          return true;
        } else {
          throw new Error('Ingresá un DNI válido');
        }
      },
    },
  },
  cvuAlias: {
    in: ['body'],
    custom: {
      options: (value, { req }) => {
        if (req.body.paymentMethodID !== PaymentMethodID.Transferencia) {
          return true;
        } else if (value && validator.isLength(value, { min: 4, max: 40 })) {
          return true;
        } else {
          throw new Error('Ingresá un CVU/CBU/Alias válido');
        }
      },
    },
  },
  paymentAccount: {
    in: ['body'],
    custom: {
      options: (value, { req }) => {
        if (req.body.paymentMethodID === PaymentMethodID.Transferencia) {
          return true;
        } else if (value && validator.isLength(value, { min: 4, max: 40 })) {
          return true;
        } else {
          throw new Error('Ingresá una cuenta de pago válida');
        }
      },
    },
  },
};

// Routes
router.post('/', verifyToken, checkSchema(addOrderSchema), addOrder);
router.get('/', verifyToken, getOrders);
router.get('/all', verifyToken, isAdmin, getAllOrders);
router.get('/:id', verifyToken, getOrder);

export default router;
