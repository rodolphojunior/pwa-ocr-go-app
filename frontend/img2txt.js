const playFlowSound = () => {
  const audio = new Audio("/sounds/flow.mp3");
  audio.play().catch((e) => console.warn("🔇 Som bloqueado pelo navegador:", e));
};

const playSuccessSound = () => {
  const audio = new Audio("/sounds/success.mp3");
  audio.play().catch((e) => console.warn("🔇 Som bloqueado pelo navegador:", e));
};

const playErrorSound = () => {
  const audio = new Audio("/sounds/error.mp3");
  audio.play().catch((e) => console.warn("🔇 Som bloqueado pelo navegador:", e));
};

document.addEventListener("DOMContentLoaded", () => {
  const modoSelect = document.getElementById("modo");
  const uploadArea = document.getElementById("upload-area");
  const cameraArea = document.getElementById("camera-area");
  const video = document.getElementById("video");
  const canvas = document.getElementById("canvas");
  const preview = document.getElementById("preview");
  const feedback = document.getElementById("feedback");
  const uploadForm = document.getElementById("upload-form");
  const btnEnviar = document.getElementById("btn-enviar");

  const token = localStorage.getItem("jwt");

  if (!token) {
    window.location.href = "/";
    return;
  }

  if (modoSelect) {
    modoSelect.addEventListener("change", () => {
      if (modoSelect.value === "camera") {
        uploadArea.style.display = "none";
        cameraArea.style.display = "block";
        navigator.mediaDevices.getUserMedia({ video: true })
          .then(stream => video.srcObject = stream)
          .catch(err => console.error("Erro ao acessar câmera:", err));
      } else {
        uploadArea.style.display = "block";
        cameraArea.style.display = "none";
      }
    });
  }

  const captureBtn = document.getElementById("capture");
  if (captureBtn) {
    captureBtn.addEventListener("click", () => {
      canvas.width = video.videoWidth;
      canvas.height = video.videoHeight;
      canvas.getContext("2d").drawImage(video, 0, 0);
      canvas.toBlob(blob => {
        const file = new File([blob], "nota.jpg", { type: "image/jpeg" });
        previewFoto(file);
        const dt = new DataTransfer();
        dt.items.add(file);
        document.getElementById("nota-fiscal").files = dt.files;
      }, "image/jpeg");
    });
  }

  if (uploadForm) {
    uploadForm.addEventListener("submit", async (e) => {
      e.preventDefault();
      const file = document.getElementById("nota-fiscal").files[0];
      if (!file) return alert("Por favor selecione ou capture uma imagem.");

      const formData = new FormData();
      formData.append("imagem", file);

      try {
        feedback.innerHTML = "⏳ Processando nota...";
        btnEnviar.disabled = true;
        btnEnviar.classList.add("loading");
        btnEnviar.innerHTML = '<span class="spinner"></span> Enviando...';

        const res = await fetch("/upload", {
          method: "POST",
          headers: {
            Authorization: `Bearer ${token}`
          },
          body: formData
        });

        if (!res.ok) throw new Error("Erro ao enviar a nota fiscal");
        const result = await res.json();
        feedback.textContent = "✅ Nota processada com sucesso!";
        playFlowSound();
        carregarNotas();
      } catch (err) {
        feedback.textContent = `❌ ${err.message}`;
        playErrorSound();

      } finally {
        btnEnviar.disabled = false;
        btnEnviar.classList.remove("loading");
        btnEnviar.innerHTML = '📤 Enviar Nota';
      }
    });
  }

  const btnCarregar = document.getElementById("btn-carregar");
  if (btnCarregar) {
    btnCarregar.addEventListener("click", carregarNotas);
  }

  const btnDeletar = document.getElementById("btn-deletar");
  if (btnDeletar) {
    btnDeletar.addEventListener("click", async () => {
      if (!confirm("Tem certeza que deseja remover todas as notas fiscais?")) return;

      try {
        const res = await fetch("/notas", {
          method: "DELETE",
          headers: {
            Authorization: `Bearer ${token}`
          }
        });
        if (!res.ok) throw new Error("Erro ao deletar notas");
        feedback.textContent = "🗑 Notas deletadas com sucesso.";
        carregarNotas();
      } catch (err) {
        feedback.textContent = `❌ ${err.message}`;
        playErrorSound();
      }
    });
  }

  async function carregarNotas() {
    try {
      const res = await fetch("/notas", {
        headers: {
          Authorization: `Bearer ${token}`
        }
      });
      if (!res.ok) throw new Error("Erro ao carregar notas");
      const notas = await res.json();
   
      const container = document.getElementById("notas-container");
      if (!notas.length) return (container.innerHTML = "<p>Nenhuma nota fiscal encontrada.</p>");

      container.innerHTML = notas.map(n => `
        <article>
          <h3>💼 ${n.empresa || "Empresa não informada"}</h3>
          <p><strong>CNPJ:</strong> ${n.cnpj || "(não informado)"}</p>
          <p><strong>Endereço:</strong> ${n.endereco || "(não informado)"}</p>
          <p><strong>Data de Emissão:</strong> ${n.data_emissao || "(não encontrada)"}</p>
          <p><strong>Valor Total:</strong> R$ ${n.valor_total.toFixed(2)}</p>
          <details>
            <summary>📋 Itens (${n.itens.length})</summary>
            <ul>
              ${n.itens.map(item => `
                <li>
                  ${item.descricao} - Qtd: ${item.quantidade}, Unit: R$ ${item.valor_unitario.toFixed(2)}
                </li>
              `).join("")}
            </ul>
          </details>
        </article>`).join("");
    } catch (err) {
      document.getElementById("notas-container").innerHTML = `<p class="error">❌ Erro ao carregar notas: ${err.message}</p>`;
      playErrorSound();
    }
  }

  function previewFoto(fileInput) {
    const file = fileInput.files ? fileInput.files[0] : fileInput;
    if (!file) return;
    const reader = new FileReader();
    reader.onload = e => {
      preview.innerHTML = `<img src="${e.target.result}" alt="Pré-visualização" style="max-width:100%; border: 2px solid #ccc; padding: 5px;">`;
    };
    reader.readAsDataURL(file);
  }
});

    if ('serviceWorker' in navigator) {
      navigator.serviceWorker.register('/service-worker.js')
        .then(reg => console.log('✅ Service Worker registrado:', reg))
        .catch(err => console.warn('⚠️ Falha no registro do Service Worker:', err));
    }


