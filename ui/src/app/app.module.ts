import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { WheelComponent } from './wheel/wheel.component';
import { StoreModule } from '@ngrx/store';
import { EffectsModule } from '@ngrx/effects';
import { playerReducer } from './store/player.reducer';

@NgModule({
  declarations: [AppComponent, WheelComponent],
  imports: [
    BrowserModule,
    AppRoutingModule,
    StoreModule.forRoot({ players: playerReducer }),
    EffectsModule.forRoot([]),
  ],
  providers: [],
  bootstrap: [AppComponent],
})
export class AppModule {}
