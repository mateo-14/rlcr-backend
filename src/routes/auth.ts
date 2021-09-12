import { Router } from 'express';
import { login, logout } from '../controllers/auth';
import isAdmin from '../middlewares/admin';
import verifyToken from '../middlewares/verifyToken';
const router = Router();

router.post('/', login);
router.post('/logout', logout);

export default router;
