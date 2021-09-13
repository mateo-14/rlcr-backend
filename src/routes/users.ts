import { Router } from 'express';
import { getUser, getAllUsers, getUserByID } from '../controllers/users';
import isAdmin from '../middlewares/admin';
import verifyToken from '../middlewares/verifyToken';
const router = Router();

router.get('/', verifyToken, getUser);

router.get('/all', verifyToken, isAdmin, getAllUsers);

router.get('/:id', verifyToken, isAdmin, getUserByID);

export default router;
