/** @format */

import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';

interface Route {
  route: string;
  icon: string;
  displayname: string;
}

@Component({
  selector: 'app-guild-admin-navbar',
  templateUrl: './guild-admin-navbar.component.html',
  styleUrls: ['./guild-admin-navbar.component.scss'],
})
export class GuildAdminNavbarComponent implements OnInit {
  public routes: Route[] = [
    {
      route: 'general',
      icon: 'settings-ol.svg',
      displayname: 'General',
    },
    {
      route: 'antiraid',
      icon: 'antiraid.svg',
      displayname: 'Antiraid',
    },
    {
      route: 'karma',
      icon: 'karma.svg',
      displayname: 'Karma',
    },
    {
      route: 'logs',
      icon: 'logs.svg',
      displayname: 'Logs',
    },
    {
      route: 'data',
      icon: 'data.svg',
      displayname: 'Data',
    },
    {
      route: 'api',
      icon: 'cloud.svg',
      displayname: 'API',
    },
  ];

  public currentPath: string;

  constructor(private router: Router, private route: ActivatedRoute) {}

  ngOnInit(): void {
    const path = window.location.pathname.split('/');
    this.currentPath = path[path.length - 1];
  }

  public navigate(r: Route) {
    const path = this.route.snapshot.url.map((u) => u.path);
    path[path.length - 1] = r.route;
    this.router.navigate(path);
  }
}
