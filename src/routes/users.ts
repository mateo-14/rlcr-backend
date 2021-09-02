import { Router } from 'express';
import { getUser } from '../controllers/users';
import verifyToken from '../middlewares/verifyToken';
const router = Router();

router.get('/', verifyToken, getUser);
export default router;
