$module('streammember', function () {
  var $static = this;
  this.$class('8933ac2f', 'SomeClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      instance.SomeInt = $t.fastbox(2, $g.________testlib.basictypes.Integer);
      return instance;
    };
    this.$typesig = function () {
      return {
      };
    };
  });

  $static.AnotherThing = function (somestream) {
    $t.streamaccess(somestream, 'SomeInt');
    return;
  };
});
