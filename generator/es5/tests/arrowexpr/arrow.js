$module('arrow', function () {
  var $static = this;
  this.$class('5df88ab1', 'SomePromise', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      return instance;
    };
    $instance.Then = function (resolve) {
      var $this = this;
      resolve($t.fastbox(true, $g.________testlib.basictypes.Boolean));
      return $this;
    };
    $instance.Catch = function (rejection) {
      var $this = this;
      return $this;
    };
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "Then|2|fd8bc7c9<69a8a979<54ff3ddf>>": true,
        "Catch|2|fd8bc7c9<69a8a979<54ff3ddf>>": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  $static.DoSomething = $t.markpromising(function (p) {
    var somebool;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      localasyncloop: while (true) {
        switch ($current) {
          case 0:
            somebool = null;
            $promise.translate(p).then(function (resolved) {
              somebool = resolved;
              $current = 1;
              $continue($resolve, $reject);
              return;
            }).catch(function (rejected) {
              $reject(rejected);
              return;
            });
            return;

          case 1:
            $resolve(somebool);
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
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      localasyncloop: while (true) {
        switch ($current) {
          case 0:
            $promise.maybe($g.arrow.DoSomething($g.arrow.SomePromise.new())).then(function ($result0) {
              $result = $result0;
              $current = 1;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 1:
            $resolve($result);
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
