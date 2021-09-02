import { Router } from 'express';
import { addOrder } from '../controllers/orders';
import verifyToken from '../middlewares/verifyToken';
const router = Router();

router.post('/', verifyToken, addOrder);

export default router;
