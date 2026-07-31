import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useForm } from 'react-hook-form';
import { useAuth } from '../context/AuthContext';
import { ApiError } from '../../../services/api';

interface SignUpForm {
  name: string;
  email: string;
  password: string;
  confirmPassword: string;
}

export default function SignUpPage() {
  const { register, handleSubmit, watch, formState: { isSubmitting } } = useForm<SignUpForm>();
  const navigate = useNavigate();
  const { register: registerUser } = useAuth();
  const [errorMsg, setErrorMsg] = useState<string | null>(null);

  const onSubmit = async (data: SignUpForm) => {
    if (data.password !== data.confirmPassword) {
      setErrorMsg('Passwords do not match.');
      return;
    }
    setErrorMsg(null);
    try {
      await registerUser(data.name, data.email, data.password);
      navigate('/dashboard', { replace: true });
    } catch (err) {
      if (err instanceof ApiError) {
        if (err.status === 409) {
          setErrorMsg('This email is already registered. Please sign in.');
        } else {
          setErrorMsg(err.message || 'Registration failed. Please try again.');
        }
      } else {
        setErrorMsg('An unexpected error occurred. Please try again.');
      }
    }
  };

  return (
    <div className="animate-in fade-in slide-in-from-bottom-4 duration-500">
      <h1 className="text-3xl font-display mb-2">Create Workspace.</h1>
      <p className="text-white/60 mb-8">Register your organization on EPMP.</p>

      {errorMsg && (
        <div className="mb-5 px-4 py-3 rounded-xl bg-red-500/10 border border-red-500/20 text-red-400 text-sm">
          {errorMsg}
        </div>
      )}

      <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col gap-5">
        <div className="flex flex-col gap-2">
          <label className="text-sm text-white/80">Full Name</label>
          <input
            {...register('name', { required: true, minLength: 2 })}
            type="text"
            placeholder="John Administrator"
            className="bg-white/5 border border-white/10 rounded-xl px-4 py-3 text-white placeholder:text-white/30 focus:outline-none focus:border-orange focus:ring-1 focus:ring-orange transition-all"
            required
          />
        </div>

        <div className="flex flex-col gap-2">
          <label className="text-sm text-white/80">Admin Email</label>
          <input
            {...register('email', { required: true })}
            type="email"
            placeholder="admin@acme.com"
            className="bg-white/5 border border-white/10 rounded-xl px-4 py-3 text-white placeholder:text-white/30 focus:outline-none focus:border-orange focus:ring-1 focus:ring-orange transition-all"
            required
          />
        </div>

        <div className="flex flex-col gap-2">
          <label className="text-sm text-white/80">Password</label>
          <input
            {...register('password', { required: true, minLength: 8 })}
            type="password"
            placeholder="Min. 8 characters"
            className="bg-white/5 border border-white/10 rounded-xl px-4 py-3 text-white placeholder:text-white/30 focus:outline-none focus:border-orange focus:ring-1 focus:ring-orange transition-all"
            required
          />
        </div>

        <div className="flex flex-col gap-2">
          <label className="text-sm text-white/80">Confirm Password</label>
          <input
            {...register('confirmPassword', {
              required: true,
              validate: value => value === watch('password') || 'Passwords do not match',
            })}
            type="password"
            placeholder="Re-enter password"
            className="bg-white/5 border border-white/10 rounded-xl px-4 py-3 text-white placeholder:text-white/30 focus:outline-none focus:border-orange focus:ring-1 focus:ring-orange transition-all"
            required
          />
        </div>

        <button
          type="submit"
          disabled={isSubmitting}
          className="mt-4 bg-white text-black font-bold py-3 rounded-xl hover:bg-orange transition-colors disabled:opacity-60 disabled:cursor-not-allowed flex items-center justify-center gap-2"
        >
          {isSubmitting ? (
            <>
              <span className="w-4 h-4 rounded-full border-2 border-black/40 border-t-transparent animate-spin" />
              Creating workspace…
            </>
          ) : 'Register Organization'}
        </button>
      </form>

      <p className="text-center text-sm text-white/60 mt-8">
        Already registered?{' '}
        <Link to="/auth/signin" className="text-white hover:text-orange font-semibold transition-colors">
          Sign In
        </Link>
      </p>
    </div>
  );
}
