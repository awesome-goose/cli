import { Component, OnInit, signal } from '@angular/core';

interface VersionInfo {
  name: string;
  version: string;
}

@Component({
  selector: 'app-root',
  standalone: true,
  template: `
    <main class="app">
      <h1>Goose + Angular</h1>
      <p>
        This page is served by the Goose SPA platform. The backend answers
        JSON under <code>/api</code>.
      </p>
      @if (info(); as i) {
        <p class="status">
          Backend says: <strong>{{ i.name }}</strong> v{{ i.version }}
        </p>
      }
      @if (error(); as e) {
        <p class="error">Backend unreachable: {{ e }}</p>
      }
    </main>
  `,
})
export class AppComponent implements OnInit {
  readonly info = signal<VersionInfo | null>(null);
  readonly error = signal<string | null>(null);

  ngOnInit(): void {
    fetch('/api/version')
      .then((res) => res.json())
      .then((data: VersionInfo) => this.info.set(data))
      .catch((err) => this.error.set(String(err)));
  }
}
