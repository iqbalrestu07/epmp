import { useState } from 'react';
import { Link, useNavigate, useLocation } from 'react-router-dom';
import { useForm } from 'react-hook-form';
import { useAuth } from '../context/AuthContext';
import { ApiError } from '../../../services/api';

interface SignInForm {
  email: string;
  password: string;
}

export default function SignInPage() {
  const { register, handleSubmit, formState: { isSubmitting } } = useForm<SignInForm>();
  const navigate = useNavigate();
  const location = useLocation();
  const { login } = useAuth();
  const [errorMsg, setErrorMsg] = useState<string | null>(null);

  const from = (location.state as { from?: { pathname: string } })?.from?.pathname || '/dashboard';

  const onSubmit = async (data: SignInForm) => {
    setErrorMsg(null);
    try {
      await login(data.email, data.password);
      navigate(from, { replace: true });
    } catch (err) {
      if (err instanceof ApiError) {
        setErrorMsg(err.message || 'Invalid email or password.');
      } else {
        setErrorMsg('An unexpected error occurred. Please try again.');
      }
    }
  };

  return (
    <div className="animate-in fade-in slide-in-from-bottom-4 duration-500">
      <h1 className="text-3xl font-display mb-2">Welcome back.</h1>
      <p className="text-white/60 mb-8">Sign in to your EPMP enterprise account.</p>

      {errorMsg && (
        <div className="mb-5 px-4 py-3 rounded-xl bg-red-500/10 border border-red-500/20 text-red-400 text-sm">
          {errorMsg}
        </div>
      )}

      <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col gap-5">
        <div className="flex flex-col gap-2">
          <label className="text-sm text-white/80">Email Address</label>
          <input
            {...register('email', { required: true })}
            type="email"
            placeholder="admin@enterprise.com"
            className="bg-white/5 border border-white/10 rounded-xl px-4 py-3 text-white placeholder:text-white/30 focus:outline-none focus:border-orange focus:ring-1 focus:ring-orange transition-all"
            required
          />
        </div>

        <div className="flex flex-col gap-2">
          <div className="flex justify-between items-center">
            <label className="text-sm text-white/80">Password</label>
            <Link to="#" className="text-xs text-orange hover:underline">Forgot password?</Link>
          </div>
          <input
            {...register('password', { required: true })}
            type="password"
            placeholder="••••••••"
            className="bg-white/5 border border-white/10 rounded-xl px-4 py-3 text-white placeholder:text-white/30 focus:outline-none focus:border-orange focus:ring-1 focus:ring-orange transition-all"
            required
          />
        </div>

        <button
          type="submit"
          disabled={isSubmitting}
          className="mt-4 bg-orange text-black font-bold py-3 rounded-xl hover:bg-white transition-colors disabled:opacity-60 disabled:cursor-not-allowed flex items-center justify-center gap-2"
        >
          {isSubmitting ? (
            <>
              <span className="w-4 h-4 rounded-full border-2 border-black/40 border-t-transparent animate-spin" />
              Signing in…
            </>
          ) : 'Sign In'}
        </button>
      </form>

      <p className="text-center text-sm text-white/60 mt-8">
        Don't have an organization account?{' '}
        <Link to="/auth/signup" className="text-white hover:text-orange font-semibold transition-colors">
          Register Organization
        </Link>
      </p>
    </div>
  );
}
