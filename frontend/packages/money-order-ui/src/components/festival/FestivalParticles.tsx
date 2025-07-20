/*
Copyright 2024 DeshChain Foundation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

import React, { useEffect, useRef } from 'react';
import { Box } from '@mui/material';
import { Festival, ParticleConfig } from '../../themes/festivals';

interface FestivalParticlesProps {
  festival: Festival;
  density?: number;
  interactive?: boolean;
}

interface Particle {
  x: number;
  y: number;
  size: number;
  speed: number;
  angle: number;
  color: string;
  type: string;
  rotation: number;
  rotationSpeed: number;
  opacity: number;
  fadeDirection: number;
}

export const FestivalParticles: React.FC<FestivalParticlesProps> = ({
  festival,
  density = 1,
  interactive = true
}) => {
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const particlesRef = useRef<Particle[]>([]);
  const animationRef = useRef<number>();
  const mouseRef = useRef({ x: 0, y: 0 });

  // Initialize particles
  const createParticle = (config: ParticleConfig): Particle => {
    const canvas = canvasRef.current;
    if (!canvas) return {} as Particle;

    const size = Math.random() * (config.size.max - config.size.min) + config.size.min;
    const speed = Math.random() * (config.speed.max - config.speed.min) + config.speed.min;
    
    return {
      x: Math.random() * canvas.width,
      y: Math.random() * canvas.height,
      size,
      speed,
      angle: Math.random() * Math.PI * 2,
      color: config.colors[Math.floor(Math.random() * config.colors.length)],
      type: config.type,
      rotation: Math.random() * Math.PI * 2,
      rotationSpeed: (Math.random() - 0.5) * 0.05,
      opacity: 1,
      fadeDirection: 1
    };
  };

  // Draw particles
  const drawParticle = (ctx: CanvasRenderingContext2D, particle: Particle) => {
    ctx.save();
    ctx.globalAlpha = particle.opacity;
    ctx.translate(particle.x, particle.y);
    ctx.rotate(particle.rotation);

    switch (particle.type) {
      case 'diya':
        drawDiya(ctx, particle);
        break;
      case 'flower':
        drawFlower(ctx, particle);
        break;
      case 'rangoli':
        drawRangoli(ctx, particle);
        break;
      case 'firework':
        drawFirework(ctx, particle);
        break;
      case 'star':
        drawStar(ctx, particle);
        break;
      case 'custom':
        drawCustom(ctx, particle);
        break;
      default:
        drawDefault(ctx, particle);
    }

    ctx.restore();
  };

  // Draw specific particle types
  const drawDiya = (ctx: CanvasRenderingContext2D, particle: Particle) => {
    // Draw flame
    ctx.fillStyle = particle.color;
    ctx.beginPath();
    ctx.moveTo(0, -particle.size / 2);
    ctx.quadraticCurveTo(-particle.size / 4, 0, 0, particle.size / 4);
    ctx.quadraticCurveTo(particle.size / 4, 0, 0, -particle.size / 2);
    ctx.fill();
    
    // Draw diya base
    ctx.fillStyle = '#8B4513';
    ctx.beginPath();
    ctx.arc(0, particle.size / 3, particle.size / 3, 0, Math.PI);
    ctx.fill();
  };

  const drawFlower = (ctx: CanvasRenderingContext2D, particle: Particle) => {
    const petals = 5;
    ctx.fillStyle = particle.color;
    
    for (let i = 0; i < petals; i++) {
      ctx.save();
      ctx.rotate((Math.PI * 2 * i) / petals);
      ctx.beginPath();
      ctx.ellipse(0, -particle.size / 3, particle.size / 4, particle.size / 2, 0, 0, Math.PI * 2);
      ctx.fill();
      ctx.restore();
    }
    
    // Center
    ctx.fillStyle = '#FFD700';
    ctx.beginPath();
    ctx.arc(0, 0, particle.size / 5, 0, Math.PI * 2);
    ctx.fill();
  };

  const drawRangoli = (ctx: CanvasRenderingContext2D, particle: Particle) => {
    const points = 8;
    ctx.strokeStyle = particle.color;
    ctx.lineWidth = 2;
    
    for (let i = 0; i < points; i++) {
      ctx.save();
      ctx.rotate((Math.PI * 2 * i) / points);
      ctx.beginPath();
      ctx.moveTo(0, 0);
      ctx.lineTo(0, -particle.size / 2);
      ctx.stroke();
      
      ctx.beginPath();
      ctx.arc(0, -particle.size / 2, particle.size / 6, 0, Math.PI * 2);
      ctx.stroke();
      ctx.restore();
    }
  };

  const drawFirework = (ctx: CanvasRenderingContext2D, particle: Particle) => {
    const sparks = 12;
    
    for (let i = 0; i < sparks; i++) {
      ctx.save();
      ctx.rotate((Math.PI * 2 * i) / sparks);
      
      const gradient = ctx.createLinearGradient(0, 0, 0, -particle.size / 2);
      gradient.addColorStop(0, particle.color + '00');
      gradient.addColorStop(1, particle.color);
      
      ctx.strokeStyle = gradient;
      ctx.lineWidth = 2;
      ctx.beginPath();
      ctx.moveTo(0, 0);
      ctx.lineTo(0, -particle.size / 2);
      ctx.stroke();
      
      ctx.restore();
    }
  };

  const drawStar = (ctx: CanvasRenderingContext2D, particle: Particle) => {
    const spikes = 5;
    const outerRadius = particle.size / 2;
    const innerRadius = particle.size / 4;
    
    ctx.fillStyle = particle.color;
    ctx.beginPath();
    
    for (let i = 0; i < spikes * 2; i++) {
      const radius = i % 2 === 0 ? outerRadius : innerRadius;
      const angle = (Math.PI * i) / spikes;
      const x = Math.cos(angle) * radius;
      const y = Math.sin(angle) * radius;
      
      if (i === 0) {
        ctx.moveTo(x, y);
      } else {
        ctx.lineTo(x, y);
      }
    }
    
    ctx.closePath();
    ctx.fill();
  };

  const drawCustom = (ctx: CanvasRenderingContext2D, particle: Particle) => {
    // For custom particles, draw a simple colored circle
    ctx.fillStyle = particle.color;
    ctx.beginPath();
    ctx.arc(0, 0, particle.size / 2, 0, Math.PI * 2);
    ctx.fill();
  };

  const drawDefault = (ctx: CanvasRenderingContext2D, particle: Particle) => {
    ctx.fillStyle = particle.color;
    ctx.beginPath();
    ctx.arc(0, 0, particle.size / 2, 0, Math.PI * 2);
    ctx.fill();
  };

  // Update particle positions
  const updateParticle = (particle: Particle) => {
    const canvas = canvasRef.current;
    if (!canvas) return;

    // Move particle
    particle.x += Math.cos(particle.angle) * particle.speed;
    particle.y += Math.sin(particle.angle) * particle.speed;
    
    // Rotate
    particle.rotation += particle.rotationSpeed;
    
    // Fade in/out
    if (particle.type === 'firework') {
      particle.opacity -= 0.01;
      if (particle.opacity <= 0) {
        // Reset firework
        particle.opacity = 1;
        particle.x = Math.random() * canvas.width;
        particle.y = Math.random() * canvas.height;
      }
    }
    
    // Wrap around edges
    if (particle.x < -particle.size) particle.x = canvas.width + particle.size;
    if (particle.x > canvas.width + particle.size) particle.x = -particle.size;
    if (particle.y < -particle.size) particle.y = canvas.height + particle.size;
    if (particle.y > canvas.height + particle.size) particle.y = -particle.size;
    
    // Interactive mouse effect
    if (interactive) {
      const dx = particle.x - mouseRef.current.x;
      const dy = particle.y - mouseRef.current.y;
      const distance = Math.sqrt(dx * dx + dy * dy);
      
      if (distance < 100) {
        const force = (100 - distance) / 100;
        particle.x += dx * force * 0.05;
        particle.y += dy * force * 0.05;
      }
    }
  };

  // Animation loop
  const animate = () => {
    const canvas = canvasRef.current;
    const ctx = canvas?.getContext('2d');
    if (!canvas || !ctx) return;

    ctx.clearRect(0, 0, canvas.width, canvas.height);
    
    particlesRef.current.forEach(particle => {
      updateParticle(particle);
      drawParticle(ctx, particle);
    });
    
    animationRef.current = requestAnimationFrame(animate);
  };

  // Initialize
  useEffect(() => {
    const canvas = canvasRef.current;
    if (!canvas) return;

    const handleResize = () => {
      canvas.width = window.innerWidth;
      canvas.height = window.innerHeight;
      
      // Recreate particles on resize
      particlesRef.current = [];
      festival.animations.particles.forEach(config => {
        const count = Math.floor(config.count * density);
        for (let i = 0; i < count; i++) {
          particlesRef.current.push(createParticle(config));
        }
      });
    };

    const handleMouseMove = (e: MouseEvent) => {
      mouseRef.current = { x: e.clientX, y: e.clientY };
    };

    handleResize();
    window.addEventListener('resize', handleResize);
    if (interactive) {
      window.addEventListener('mousemove', handleMouseMove);
    }
    
    animate();

    return () => {
      window.removeEventListener('resize', handleResize);
      if (interactive) {
        window.removeEventListener('mousemove', handleMouseMove);
      }
      if (animationRef.current) {
        cancelAnimationFrame(animationRef.current);
      }
    };
  }, [festival, density, interactive]);

  return (
    <Box
      component="canvas"
      ref={canvasRef}
      sx={{
        position: 'fixed',
        top: 0,
        left: 0,
        width: '100%',
        height: '100%',
        pointerEvents: 'none',
        zIndex: 1
      }}
    />
  );
};