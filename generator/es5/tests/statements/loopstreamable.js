$module('loopstreamable', function () {
  var $static = this;
  this.$class('4f99dc7f', 'SomeStreamable', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      return instance;
    };
    $instance.Stream = function () {
      var $this = this;
      return $g.loopstreamable.SomeStream.new();
    };
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "Stream|2|fd8bc7c9<1d733667<9706e8ab>>": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  this.$class('e6a7b534', 'SomeStream', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      instance.wasChecked = $t.fastbox(false, $g.________testlib.basictypes.Boolean);
      return instance;
    };
    $instance.Next = function () {
      var $this = this;
      var r;
      r = $this.wasChecked;
      $this.wasChecked = $t.fastbox(true, $g.________testlib.basictypes.Boolean);
      return $g.________testlib.basictypes.Tuple($g.________testlib.basictypes.Boolean, $g.________testlib.basictypes.Boolean).Build($t.fastbox(true, $g.________testlib.basictypes.Boolean), $t.fastbox(!r.$wrapped, $g.________testlib.basictypes.Boolean));
    };
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "Next|2|fd8bc7c9<58998129<9706e8ab,9706e8ab>>": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  $static.DoSomething = $t.markpromising(function (somethingElse) {
    var $result;
    var $temp0;
    var $temp1;
    var something;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      localasyncloop: while (true) {
        switch ($current) {
          case 0:
            $t.fastbox(1234, $g.________testlib.basictypes.Integer);
            $current = 1;
            continue localasyncloop;

          case 1:
            $temp1 = somethingElse.Stream();
            $current = 2;
            continue localasyncloop;

          case 2:
            $promise.maybe($temp1.Next()).then(function ($result0) {
              $temp0 = $result0;
              $result = $temp0;
              $current = 3;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 3:
            something = $temp0.First;
            if ($temp0.Second.$wrapped) {
              $current = 4;
              continue localasyncloop;
            } else {
              $current = 5;
              continue localasyncloop;
            }
            break;

          case 4:
            $t.fastbox(7654, $g.________testlib.basictypes.Integer);
            $current = 2;
            continue localasyncloop;

          case 5:
            $t.fastbox(5678, $g.________testlib.basictypes.Integer);
            $resolve();
            return;

          default:
            $resolve();
            return;
        }
      }
    };
    return $promise.new($continue);
  });
  $static.TEST = $t.markpromising(function () {
    var $result;
    var $temp0;
    var $temp1;
    var i;
    var result;
    var s;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      localasyncloop: while (true) {
        switch ($current) {
          case 0:
            result = $t.fastbox('noloop', $g.________testlib.basictypes.String);
            s = $g.loopstreamable.SomeStreamable.new();
            $current = 1;
            continue localasyncloop;

          case 1:
            $temp1 = s.Stream();
            $current = 2;
            continue localasyncloop;

          case 2:
            $promise.maybe($temp1.Next()).then(function ($result0) {
              $temp0 = $result0;
              $result = $temp0;
              $current = 3;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 3:
            i = $temp0.First;
            if ($temp0.Second.$wrapped) {
              $current = 4;
              continue localasyncloop;
            } else {
              $current = 5;
              continue localasyncloop;
            }
            break;

          case 4:
            result = i;
            $current = 2;
            continue localasyncloop;

          case 5:
            $resolve(result);
            return;

          default:
            $resolve();
            return;
        }
      }
    };
    return $promise.new($continue);
  });
});
