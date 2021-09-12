import { Router } from 'express';
import { getUser, getAllUsers } from '../controllers/users';
import isAdmin from '../middlewares/admin';
import verifyToken from '../middlewares/verifyToken';
const router = Router();

router.get('/', verifyToken, getUser);

router.get('/all', verifyToken, isAdmin, getAllUsers);

export default router;
